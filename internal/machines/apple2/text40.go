package apple2

import (
	"basics/internal/video"
	ebitenrenderer "basics/internal/video/ebiten"
	"basics/internal/video/text"
	"bufio"
	"image/color"
	"io"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Interface locale (MINIMALE) pour éviter de polluer video.Renderer
type ebitenRenderer interface {
	Draw(screen *ebiten.Image)
	Layout(w, h int) (int, int)
}

type Text40 struct {
	Mode     *text.TextMode
	renderer video.Renderer

	in  *bufio.Reader
	out io.Writer

	// For INPUT
	inputBuffer []rune
	lineReady   bool

	// For GET
	getActive bool
	getChan   chan rune

	// Blinking cursor
	cursorVisible bool
	blinkCounter  int
	inInput       bool

	// Input is allowed
	allowInput bool
}

func NewText40(renderer video.Renderer) *Text40 {
	mode := text.NewTextMode(
		renderer,
		40, 24, // Apple II Text 40 colonnes
		7, 8, // font 7x8
		1, 0, // blanc sur noir
	)

	return &Text40{
		Mode:        mode,
		renderer:    renderer,
		in:          bufio.NewReader(strings.NewReader("")),
		out:         io.Discard,
		inputBuffer: make([]rune, 0, 64),
		lineReady:   false,
		allowInput:  false,
	}
}

// --------------------
// video.Device
// --------------------

func (t *Text40) Clear() {
	t.Mode.Home()
}

func (t *Text40) PrintChar(r rune) {
	t.Mode.PutChar(r)
}

func (t *Text40) PrintString(s string) {
	t.Mode.Print(s)
}

func (t *Text40) SetCursorX(x int) {
	t.Mode.SetCursor(x, t.Mode.CursorY())
}

func (t *Text40) SetCursorY(y int) {
	t.Mode.SetCursor(t.Mode.CursorX(), y)
}

func (t *Text40) Plot(x, y int) {
	// ignoré en mode texte
}

func (t *Text40) Render() {
	t.Mode.Render()
}

// --------------------
// I/O
// --------------------

func (t *Text40) SetInput(r io.Reader) {
	t.in = bufio.NewReader(r)
}

func (t *Text40) SetOutput(w io.Writer) {
	t.out = w
}

// --------------------
// Ebiten integration
// --------------------

func (t *Text40) Update() error {
	if !t.inInput {
		t.cursorVisible = false
		t.blinkCounter = 0
		return nil
	}

	t.blinkCounter++
	if t.blinkCounter >= 30 { // ~0.5s à 60 FPS
		t.cursorVisible = !t.cursorVisible
		t.blinkCounter = 0
	}

	return nil
}

func (t *Text40) Draw(screen *ebiten.Image) {
	// Gestion du curseur clignotant
	if t.inInput && t.cursorVisible {
		t.Mode.PutChar('░')
		t.SetCursorX(t.Mode.CursorX() - 1)
	} else if t.inInput && !t.cursorVisible {
		t.Mode.PutChar(' ')
		t.SetCursorX(t.Mode.CursorX() - 1)
	}

	// Demande au TextMode de rasteriser le buffer
	t.Mode.Render()
	t.Mode.Renderer.(*ebitenrenderer.Renderer).BlitTo(screen)

	// Demande au renderer Ebiten d’afficher l’image
	if r, ok := t.renderer.(ebitenRenderer); ok {
		screen.Fill(color.Black)
		r.Draw(screen)
	}
}

func (t *Text40) Layout(w, h int) (int, int) {
	if r, ok := t.renderer.(ebitenRenderer); ok {
		return r.Layout(w, h)
	}
	return w, h
}

// --------------------
// Input & cursor movement
// --------------------
func (t *Text40) ReadLine() (string, error) {
	t.BeginInput()
	defer t.EndInput()

	for !t.lineReady {
		// attente active mais NON bloquante
		time.Sleep(5 * time.Millisecond)
	}

	line := string(t.inputBuffer)

	t.inputBuffer = t.inputBuffer[:0]
	t.lineReady = false

	// comportement AppleSoft : retour à la ligne automatique
	t.Mode.NewLine()

	return line, nil
}

func (t *Text40) InputRune(r rune) {
	if !t.allowInput {
		return
	}

	t.eraseCursorIfVisible()

	t.inputBuffer = append(t.inputBuffer, r)
	t.Mode.PutChar(r)
}

func (t *Text40) Backspace() {
	if !t.allowInput || len(t.inputBuffer) == 0 {
		return
	}

	t.eraseCursorIfVisible()
	t.Mode.Backspace()

	t.inputBuffer = t.inputBuffer[:len(t.inputBuffer)-1]
}

func (t *Text40) Enter() {
	if !t.allowInput {
		return
	}

	t.eraseCursorIfVisible()
	t.EndInput()
	t.lineReady = true
}

func (t *Text40) BeginInput() {
	t.inInput = true
	t.allowInput = true
	t.cursorVisible = true
	t.blinkCounter = 0
}

func (t *Text40) EndInput() {
	t.eraseCursorIfVisible()
	t.inInput = false
	t.allowInput = false
	t.cursorVisible = false
}

func (t *Text40) BeginGet() {
	t.getActive = true
	t.getChan = make(chan rune, 1)
}

func (t *Text40) EndGet() {
	t.getActive = false
}

func (t *Text40) PushGetRune(r rune) {
	if t.getActive {
		select {
		case t.getChan <- r:
		default:
		}
	}
}

func (t *Text40) GetChar() (rune, error) {
	t.BeginGet()
	r := <-t.getChan
	t.EndGet()
	return r, nil
}

func (t *Text40) IsGetActive() bool {
	return t.getActive
}

func (t *Text40) eraseCursorIfVisible() {
	if t.inInput && t.cursorVisible {
		// remplacer le curseur par un espace
		t.Mode.PutChar(' ')
		t.SetCursorX(t.Mode.CursorX() - 1)
		t.cursorVisible = false
		t.blinkCounter = 0
	}
}

func (t *Text40) DisableKeyboard() {
	t.allowInput = false
	t.inInput = false
}
