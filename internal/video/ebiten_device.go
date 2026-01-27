package video

import "github.com/hajimehoshi/ebiten/v2"

type EbitenDevice interface {
	Device
	Draw(screen *ebiten.Image)
	Layout(w, h int) (int, int)
}
