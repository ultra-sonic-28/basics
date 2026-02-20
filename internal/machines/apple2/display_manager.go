package apple2

import (
	"basics/internal/constants"
	"basics/internal/video"
	ebitenrenderer "basics/internal/video/ebiten"
	"basics/internal/video/font"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
)

type DisplayManager struct {
	machine byte
	modes   map[video.ModeID]video.Mode
	current video.Mode
}

var _ video.EbitenDevice = (*DisplayManager)(nil)

func NewAppleDisplay(
	scale int,
	defaultMode video.ModeID,
) *DisplayManager {
	basicType := constants.BASIC_APPLE

	// 40 colonnes = 280x192
	r40 := ebitenrenderer.New(
		scale*280, scale*192,
		scale,
		Palette(),
		font.DefaultFontForMode(basicType),
	)
	text40 := NewAppleText(
		r40,
		40, 24,
		7, 8,
		ModeAppleText40,
	)

	// 80 colonnes = 560x192
	r80 := ebitenrenderer.New(
		scale*560, scale*192,
		scale,
		Palette(),
		font.DefaultFontForMode(basicType),
	)
	text80 := NewAppleText(
		r80,
		80, 24,
		7, 8,
		ModeAppleText80,
	)

	modes := map[video.ModeID]video.Mode{
		ModeAppleText40: text40,
		ModeAppleText80: text80,
	}

	current, ok := modes[defaultMode]
	if !ok {
		current = text40
	}

	return &DisplayManager{
		machine: constants.BASIC_APPLE,
		modes:   modes,
		current: current,
	}
}

// ---------------------------------------------------

func (d *DisplayManager) Switch(id video.ModeID) {
	if m, ok := d.modes[id]; ok {
		if d.current != m {
			// Copier le contenu de l'ancien buffer vers le nouveau
			oldMode, okOld := d.current.(*AppleText)
			newMode, okNew := m.(*AppleText)
			if okOld && okNew {
				newMode.CopyFrom(oldMode)
			}

			d.current = m
			d.Render() // Forcer un premier rendu dans le nouveau mode
		}
	}
}

func (d *DisplayManager) Modes() []video.ModeID {
	ids := make([]video.ModeID, 0, len(d.modes))
	for id := range d.modes {
		ids = append(ids, id)
	}
	return ids
}

func (d *DisplayManager) ModeInfo(id video.ModeID) (video.ModeInfo, bool) {
	m, ok := d.modes[id]
	if !ok {
		return video.ModeInfo{}, false
	}
	return m.Info(), true
}

func (d *DisplayManager) DefaultMode() video.ModeID {
	return d.current.Info().ID
}

// ---------------------------------------------------
// Ebiten delegation
// ---------------------------------------------------

func (d *DisplayManager) Draw(screen *ebiten.Image) {
	if dev, ok := d.current.(video.EbitenDevice); ok {
		dev.Draw(screen)
	}
}

func (d *DisplayManager) Layout(w, h int) (int, int) {
	if dev, ok := d.current.(video.EbitenDevice); ok {
		return dev.Layout(w, h)
	}
	return 0, 0
}

func (d *DisplayManager) Width() int {
	if dev, ok := d.current.(video.EbitenDevice); ok {
		return dev.Width()
	}
	return 0
}

func (d *DisplayManager) Height() int {
	if dev, ok := d.current.(video.EbitenDevice); ok {
		return dev.Height()
	}
	return 0
}

func (d *DisplayManager) Scale() int {
	if dev, ok := d.current.(video.EbitenDevice); ok {
		return dev.Scale()
	}
	return 1
}

func (d *DisplayManager) Update() error {
	if m, ok := d.current.(interface{ Update() error }); ok {
		return m.Update()
	}
	return nil
}

// ---------------------------------------------------
// Device delegation
// ---------------------------------------------------

func (d *DisplayManager) device() video.Device {
	return d.current.(video.Device)
}

func (d *DisplayManager) Clear()               { d.device().Clear() }
func (d *DisplayManager) SetInverse(v bool)    { d.device().SetInverse(v) }
func (d *DisplayManager) SetFlash(v bool)      { d.device().SetFlash(v) }
func (d *DisplayManager) PrintChar(r rune)     { d.device().PrintChar(r) }
func (d *DisplayManager) PrintString(s string) { d.device().PrintString(s) }
func (d *DisplayManager) SetCursorX(x int)     { d.device().SetCursorX(x) }
func (d *DisplayManager) SetCursorY(y int)     { d.device().SetCursorY(y) }

func (d *DisplayManager) SwitchMode(slot int) {
	switch slot {
	case 0:
		d.Switch(ModeAppleText40)
	case 3:
		d.Switch(ModeAppleText80)
	}
}

func (d *DisplayManager) Plot(x, y int)             { d.device().Plot(x, y) }
func (d *DisplayManager) ReadLine() (string, error) { return d.device().ReadLine() }
func (d *DisplayManager) GetChar() (rune, error)    { return d.device().GetChar() }
func (d *DisplayManager) SetOutput(w io.Writer)     { d.device().SetOutput(w) }
func (d *DisplayManager) DisableKeyboard()          { d.device().DisableKeyboard() }
func (d *DisplayManager) Render()                   { d.device().Render() }

// ---------------------------------------------------
// Input delegation (for EbitenApp)
// ---------------------------------------------------

func (d *DisplayManager) IsGetActive() bool {
	if m, ok := d.current.(interface{ IsGetActive() bool }); ok {
		return m.IsGetActive()
	}
	return false
}

func (d *DisplayManager) PushGetRune(r rune) {
	if m, ok := d.current.(interface{ PushGetRune(rune) }); ok {
		m.PushGetRune(r)
	}
}

func (d *DisplayManager) InputRune(r rune) {
	if m, ok := d.current.(interface{ InputRune(rune) }); ok {
		m.InputRune(r)
	}
}

func (d *DisplayManager) Backspace() {
	if m, ok := d.current.(interface{ Backspace() }); ok {
		m.Backspace()
	}
}

func (d *DisplayManager) Enter() {
	if m, ok := d.current.(interface{ Enter() }); ok {
		m.Enter()
	}
}
