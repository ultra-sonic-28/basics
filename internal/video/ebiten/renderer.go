package ebitenrenderer

import (
	"basics/internal/video"
	"basics/internal/video/font"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	width, height int
	scale         int

	palette video.Palette

	framebuffer *image.RGBA
	font        *font.BitmapFont
}

func New(
	width, height int,
	scale int,
	palette video.Palette,
	font *font.BitmapFont,
) *Renderer {
	fb := image.NewRGBA(image.Rect(0, 0, width, height))

	return &Renderer{
		width:       width,
		height:      height,
		scale:       scale,
		palette:     palette,
		font:        font,
		framebuffer: fb,
	}
}

func (r *Renderer) Width() int  { return r.width }
func (r *Renderer) Height() int { return r.height }

func (r *Renderer) Clear() {
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.framebuffer.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
}

func (r *Renderer) DrawPixel(x, y int, c int) {
	if x < 0 || y < 0 || x >= r.width || y >= r.height {
		return
	}

	if c < 0 || c >= len(r.palette) {
		return
	}

	r.framebuffer.SetRGBA(x, y, r.palette[c])
}

func (r *Renderer) DrawGlyph(x, y int, glyph rune, fg, bg int) {
	bitmap := r.font.Glyph(glyph)

	for row := 0; row < r.font.Height; row++ {
		bits := bitmap[row]

		for col := 0; col < r.font.Width; col++ {
			px := x + col
			py := y + row

			if px < 0 || py < 0 || px >= r.width || py >= r.height {
				continue
			}

			mask := byte(1 << col)
			// lecture MSB -> LSB
			//shift := (r.font.Width - 1) - col
			//mask := byte(1 << shift)
			if bits&mask != 0 {
				r.framebuffer.SetRGBA(px, py, r.palette[fg])
			} else {
				r.framebuffer.SetRGBA(px, py, r.palette[bg])
			}
		}
	}
}

func (r *Renderer) BlitTo(screen *ebiten.Image) {
	img := ebiten.NewImageFromImage(r.framebuffer)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(r.scale), float64(r.scale))
	screen.DrawImage(img, op)
}
