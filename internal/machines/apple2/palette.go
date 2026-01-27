package apple2

import (
	"basics/internal/video"
	"image/color"
)

func Palette() video.Palette {
	return video.Palette{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // noir
		color.RGBA{0xff, 0xff, 0xff, 0xff}, // blanc
		color.RGBA{0x00, 0xff, 0x00, 0xff}, // vert
		color.RGBA{0xff, 0x80, 0x00, 0xff}, // orange
	}
}
