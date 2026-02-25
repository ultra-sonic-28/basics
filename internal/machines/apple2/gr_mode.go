package apple2

import (
	"bufio"
	"io"
	"strings"
	"time"

	"basics/internal/video"
	ebitenrenderer "basics/internal/video/ebiten"
	"basics/internal/video/text"

	"github.com/hajimehoshi/ebiten/v2"
)

type AppleGR struct {
	ModeID   video.ModeID
	renderer video.Renderer
	isFull   bool

	// 40x48 blocks grid (0-15 color index)
	Blocks [40][48]int

	// Text part for mixed mode (bottom 4 lines)
	TextPart *text.TextMode

	// Current graphics color
	CurrentColor int

	// Device common fields
	in  *bufio.Reader
	out io.Writer

	// For INPUT (in mixed mode)
	inputBuffer []rune
	lineReady   bool
	getActive   bool
	getChan     chan rune
	allowInput  bool
	inInput     bool

	// Cursor for mixed mode text
	cursorVisible bool
	blinkCounter  int
}

var _ video.EbitenDevice = (*AppleGR)(nil)

func NewAppleGR(
	renderer video.Renderer,
	modeID video.ModeID,
	isFull bool,
) *AppleGR {
	// Mixed mode has 4 text lines at bottom
	// Text resolution is always 40 cols for Apple II GR mixed mode
	textPart := text.NewTextMode(
		renderer,
		40, 24, // Full screen buffer but we only show bottom 4 lines
		7, 8,
		15, 0, // White on black
	)
	// Position cursor at line 20 (first of the 4 bottom lines)
	textPart.SetCursor(0, 20)

	return &AppleGR{
		ModeID:       modeID,
		renderer:     renderer,
		isFull:       isFull,
		TextPart:     textPart,
		CurrentColor: 0,
		in:           bufio.NewReader(strings.NewReader("")),
		out:          io.Discard,
		inputBuffer:  make([]rune, 0, 64),
	}
}

func (g *AppleGR) Info() video.ModeInfo {
	return video.ModeInfo{
		ID:     g.ModeID,
		Name:   "Apple II GR",
		Width:  280,
		Height: 192,
		Text:   false,
	}
}

func (g *AppleGR) Reset() {
	g.Clear()
}

func (g *AppleGR) Update() error {
	if !g.inInput {
		g.cursorVisible = false
		g.blinkCounter = 0
		return nil
	}

	g.blinkCounter++
	if g.blinkCounter >= 30 {
		g.cursorVisible = !g.cursorVisible
		g.blinkCounter = 0
	}
	return nil
}

func (g *AppleGR) Draw(screen *ebiten.Image) {
	if r, ok := g.renderer.(*ebitenrenderer.Renderer); ok {
		r.NextFrame()
	}

	// 1. Draw Blocks
	// Blocks are 7x4 pixels
	rows := 48
	if !g.isFull {
		rows = 40
	}

	for y := 0; y < rows; y++ {
		for x := 0; x < 40; x++ {
			c := g.Blocks[x][y]
			// A block is 7x4 pixels
			for py := 0; py < 4; py++ {
				for px := 0; px < 7; px++ {
					g.renderer.DrawPixel(x*7+px, y*4+py, c)
				}
			}
		}
	}

	// 2. Draw Text (if mixed)
	if !g.isFull {
		// Draw only bottom 4 lines (rows 20 to 23 of the 24-row buffer)
		for ty := 20; ty < 24; ty++ {
			for tx := 0; tx < 40; tx++ {
				cell := g.TextPart.Buffer.CellAt(tx, ty)
				// ty*8 gives pixel Y position
				g.renderer.DrawGlyph(tx*7, ty*8, cell.Glyph, cell.FG, cell.BG)
			}
		}

		// Cursor
		if g.inInput && g.cursorVisible {
			cx, cy := g.TextPart.CursorX(), g.TextPart.CursorY()
			if cy >= 20 {
				g.renderer.DrawGlyph(cx*7, cy*8, '░', 15, 0)
			}
		}
	}

	if r, ok := g.renderer.(interface {
		BlitTo(screen *ebiten.Image)
	}); ok {
		r.BlitTo(screen)
	}
}

func (g *AppleGR) Layout(w, h int) (int, int) {
	return g.renderer.Width(), g.renderer.Height()
}

func (g *AppleGR) Width() int  { return g.renderer.Width() }
func (g *AppleGR) Height() int { return g.renderer.Height() }
func (g *AppleGR) Scale() int  { return g.renderer.Scale() }

// Device API
func (g *AppleGR) Clear() {
	for x := 0; x < 40; x++ {
		for y := 0; y < 48; y++ {
			g.Blocks[x][y] = 0
		}
	}
	if !g.isFull {
		g.TextPart.Home()
		g.TextPart.SetCursor(0, 20)
	}
}

func (g *AppleGR) SetInverse(v bool) { g.TextPart.SetInverse(v) }
func (g *AppleGR) SetFlash(v bool)   { g.TextPart.SetFlash(v) }

func (g *AppleGR) PrintChar(r rune) {
	if !g.isFull {
		g.TextPart.PutChar(r)
		// Ensure cursor stays in bottom 4 lines?
		// Applesoft scrolls the bottom 4 lines if needed.
		if g.TextPart.CursorY() < 20 {
			g.TextPart.SetCursor(g.TextPart.CursorX(), 20)
		}
	}
}

func (g *AppleGR) PrintString(s string) {
	for _, r := range s {
		g.PrintChar(r)
	}
}

func (g *AppleGR) SetCursorX(x int) { g.TextPart.SetCursor(x, g.TextPart.CursorY()) }
func (g *AppleGR) SetCursorY(y int) { g.TextPart.SetCursor(g.TextPart.CursorX(), y) }
func (g *AppleGR) CursorX() int     { return g.TextPart.CursorX() }
func (g *AppleGR) CursorY() int     { return g.TextPart.CursorY() }

func (g *AppleGR) SwitchMode(slot int) {} // Managed by DisplayManager

func (g *AppleGR) Plot(x, y int) {
	if x >= 0 && x < 40 && y >= 0 && y < 48 {
		g.Blocks[x][y] = g.CurrentColor
	}
}

func (g *AppleGR) SetColor(c int) {
	if c >= 0 && c <= 15 {
		g.CurrentColor = c
	}
}

func (g *AppleGR) ReadLine() (string, error) {
	if g.isFull {
		return "", nil // Text input not supported in full GR
	}
	g.BeginInput()
	defer g.EndInput()
	for !g.lineReady {
		time.Sleep(5 * time.Millisecond)
	}
	line := string(g.inputBuffer)
	g.inputBuffer = g.inputBuffer[:0]
	g.lineReady = false
	g.TextPart.NewLine()
	if g.TextPart.CursorY() < 20 {
		g.TextPart.SetCursor(g.TextPart.CursorX(), 20)
	}
	return line, nil
}

func (g *AppleGR) GetChar() (rune, error) {
	g.BeginGet()
	r := <-g.getChan
	g.EndGet()
	return r, nil
}

func (g *AppleGR) SetOutput(w io.Writer) { g.out = w }
func (g *AppleGR) DisableKeyboard()      { g.allowInput = false; g.inInput = false }
func (g *AppleGR) Render()               {} // Handled in Draw

// Input Handling
func (g *AppleGR) InputRune(r rune) {
	if !g.allowInput {
		return
	}
	g.inputBuffer = append(g.inputBuffer, r)
	g.TextPart.PutChar(r)
}

func (g *AppleGR) Backspace() {
	if !g.allowInput || len(g.inputBuffer) == 0 {
		return
	}
	g.TextPart.Backspace()
	g.inputBuffer = g.inputBuffer[:len(g.inputBuffer)-1]
}

func (g *AppleGR) Enter() {
	if !g.allowInput {
		return
	}
	g.EndInput()
	g.lineReady = true
}

func (g *AppleGR) IsGetActive() bool { return g.getActive }
func (g *AppleGR) PushGetRune(r rune) {
	if g.getActive {
		select {
		case g.getChan <- r:
		default:
		}
	}
}

func (g *AppleGR) BeginInput() {
	g.inInput = true
	g.allowInput = true
	g.cursorVisible = true
	g.blinkCounter = 0
}

func (g *AppleGR) EndInput() {
	g.inInput = false
	g.allowInput = false
	g.cursorVisible = false
}

func (g *AppleGR) BeginGet() {
	g.getActive = true
	g.getChan = make(chan rune, 1)
}

func (g *AppleGR) EndGet() { g.getActive = false }
