package apple2

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"basics/internal/video"
	ebitenrenderer "basics/internal/video/ebiten"
	"basics/internal/video/text"

	"github.com/hajimehoshi/ebiten/v2"
)

type AppleText struct {
	Mode     *text.TextMode
	renderer video.Renderer

	cols     int
	rows     int
	charW    int
	charH    int
	logicalW int
	logicalH int
	modeID   video.ModeID

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

	// Printing in normal or inverse mode
	inverse bool

	// Printed chars are flashing
	flash bool
}

var _ video.EbitenDevice = (*AppleText)(nil)

func NewAppleText(
	renderer video.Renderer,
	cols int,
	rows int,
	charW int,
	charH int,
	modeID video.ModeID,
) *AppleText {

	logicalW := cols * charW
	logicalH := rows * charH

	mode := text.NewTextMode(
		renderer,
		cols, rows,
		7, 8, // font 7x8
		15, 0, // blanc sur noir
	)
	return &AppleText{
		Mode:        mode,
		renderer:    renderer,
		cols:        cols,
		rows:        rows,
		charW:       charW,
		charH:       charH,
		logicalW:    logicalW,
		logicalH:    logicalH,
		modeID:      modeID,
		in:          bufio.NewReader(strings.NewReader("")),
		out:         io.Discard,
		inputBuffer: make([]rune, 0, 64),
		lineReady:   false,
		allowInput:  false,
		inverse:     false,
		flash:       false,
	}
}

// ---------------------------------------------------
// Mode
// ---------------------------------------------------

func (t *AppleText) Info() video.ModeInfo {
	return video.ModeInfo{
		ID:     t.modeID,
		Name:   fmt.Sprintf("Apple II Text %d", t.cols),
		Width:  t.logicalW,
		Height: t.logicalH,
		Text:   true,
	}
}

func (t *AppleText) Reset() {
	t.Clear()
}

// ---------------------------------------------------
// EbitenDevice
// ---------------------------------------------------

func (t *AppleText) Update() error {
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

func (t *AppleText) Draw(screen *ebiten.Image) {
	scale := t.Scale()
	offset := float64(GetMonitorOffset(scale))

	// 1. Draw Monitor Frame
	DrawMonitorFrame(screen, t.Width(), t.Height(), scale)

	// Gestion du curseur clignotant
	if t.inInput && t.cursorVisible {
		t.Mode.PutChar('░')
		t.SetCursorX(t.Mode.CursorX() - 1)
	} else if t.inInput && !t.cursorVisible {
		t.Mode.PutChar(' ')
		t.SetCursorX(t.Mode.CursorX() - 1)
	}

	if r, ok := t.renderer.(*ebitenrenderer.Renderer); ok {
		r.NextFrame()
	}

	// Demande au TextMode de rasteriser le buffer
	t.Mode.Render()

	// 2. Render content with offset
	if r, ok := t.renderer.(interface {
		BlitTo(screen *ebiten.Image, x, y float64)
	}); ok {
		r.BlitTo(screen, offset, offset)
	}
}

func (t *AppleText) Layout(w, h int) (int, int) {
	return t.Width(), t.Height()
}

func (t *AppleText) Width() int {
	return t.renderer.Width() + 2*GetMonitorOffset(t.Scale())
}

func (t *AppleText) Height() int {
	return t.renderer.Height() + 2*GetMonitorOffset(t.Scale())
}

func (t *AppleText) Scale() int {
	return t.renderer.Scale()
}

// ---------------------------------------------------
// Device API
// ---------------------------------------------------

func (t *AppleText) Clear() {
	t.Mode.Home()
}

func (t *AppleText) SetInverse(v bool) {
	t.inverse = v
	t.Mode.SetInverse(v)
}

func (t *AppleText) SetFlash(v bool) {
	t.flash = v
	t.Mode.SetFlash(v)
}

func (t *AppleText) PrintChar(r rune) {
	t.Mode.PutChar(r)
}

func (t *AppleText) PrintString(s string) {
	t.Mode.Print(s)
}

func (t *AppleText) SetCursorX(x int) {
	t.Mode.SetCursor(x, t.Mode.CursorY())
}

func (t *AppleText) SetCursorY(y int) {
	t.Mode.SetCursor(t.Mode.CursorX(), y)
}

func (t *AppleText) CursorX() int {
	return t.Mode.CursorX()
}

func (t *AppleText) CursorY() int {
	return t.Mode.CursorY()
}

func (t *AppleText) SwitchMode(slot int) {
	// Le switch de mode est géré par le DisplayManager qui contient les instances AppleText.
	// Cette méthode est ici pour satisfaire l'interface video.Device.
}

func (t *AppleText) Plot(x, y int) {
	// non utilisé en mode texte
}

func (t *AppleText) SetColor(c int) {
	// non utilisé en mode texte pour l'instant (utilisé pour PLOT)
}

func (t *AppleText) ReadLine() (string, error) {
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

func (t *AppleText) GetChar() (rune, error) {
	t.BeginGet()
	r := <-t.getChan
	t.EndGet()
	return r, nil
}

func (t *AppleText) SetOutput(w io.Writer) {
	t.out = w
}

func (t *AppleText) DisableKeyboard() {
	t.allowInput = false
	t.inInput = false
}

func (t *AppleText) Render() {
	t.Mode.Render()
}

// ---------------------------------------------------
// Input Handling
// ---------------------------------------------------

func (t *AppleText) InputRune(r rune) {
	if !t.allowInput {
		return
	}

	t.eraseCursorIfVisible()

	t.inputBuffer = append(t.inputBuffer, r)
	t.Mode.PutChar(r)
}

func (t *AppleText) Backspace() {
	if !t.allowInput || len(t.inputBuffer) == 0 {
		return
	}

	t.eraseCursorIfVisible()
	t.Mode.Backspace()

	t.inputBuffer = t.inputBuffer[:len(t.inputBuffer)-1]
}

func (t *AppleText) Enter() {
	if !t.allowInput {
		return
	}

	t.eraseCursorIfVisible()
	t.EndInput()
	t.lineReady = true
}

func (t *AppleText) IsGetActive() bool {
	return t.getActive
}

func (t *AppleText) PushGetRune(r rune) {
	if t.getActive {
		select {
		case t.getChan <- r:
		default:
		}
	}
}

func (t *AppleText) BeginInput() {
	t.inInput = true
	t.allowInput = true
	t.cursorVisible = true
	t.blinkCounter = 0
}

func (t *AppleText) EndInput() {
	t.eraseCursorIfVisible()
	t.inInput = false
	t.allowInput = false
	t.cursorVisible = false
}

func (t *AppleText) BeginGet() {
	t.getActive = true
	t.getChan = make(chan rune, 1)
}

func (t *AppleText) EndGet() {
	t.getActive = false
}

func (t *AppleText) eraseCursorIfVisible() {
	if t.inInput && t.cursorVisible {
		// remplacer le curseur par un espace
		t.Mode.PutChar(' ')
		t.SetCursorX(t.Mode.CursorX() - 1)
		t.cursorVisible = false
		t.blinkCounter = 0
	}
}

func (t *AppleText) CopyFrom(other *AppleText) {
	// Effacer le buffer actuel
	t.Mode.Home()

	// Copier les cellules (limité par les dimensions minimales)
	srcCols, srcRows := other.cols, other.rows
	dstCols, dstRows := t.cols, t.rows

	maxCols := srcCols
	if dstCols < maxCols {
		maxCols = dstCols
	}
	maxRows := srcRows
	if dstRows < maxRows {
		maxRows = dstRows
	}

	for y := 0; y < maxRows; y++ {
		for x := 0; x < maxCols; x++ {
			cell := other.Mode.Buffer.CellAt(x, y)
			t.Mode.Buffer.SetCell(x, y, cell.Glyph, cell.FG, cell.BG, cell.Flash)
		}
	}

	// Copier la position du curseur (en s'assurant qu'elle reste dans les limites)
	newX, newY := other.Mode.CursorX(), other.Mode.CursorY()
	if newX >= dstCols {
		newX = dstCols - 1
	}
	if newY >= dstRows {
		newY = dstRows - 1
	}
	t.Mode.SetCursor(newX, newY)
}
