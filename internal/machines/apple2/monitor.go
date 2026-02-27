package apple2

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed assets/logo.png
var logoData []byte

var (
	logoImage     *ebiten.Image
	loadLogoOnce  sync.Once
)

func getLogoImage() *ebiten.Image {
	loadLogoOnce.Do(func() {
		img, _, err := image.Decode(bytes.NewReader(logoData))
		if err != nil {
			log.Printf("failed to decode logo: %v", err)
			return
		}
		logoImage = ebiten.NewImageFromImage(img)
	})
	return logoImage
}

const (
	// MonitorPaddingBase is the base space between the BASIC screen and the monitor bezel
	MonitorPaddingBase = 8
	// MonitorBezelBase is the base thickness of the monitor frame (bezel)
	MonitorBezelBase = 10
	// MonitorInternalPaddingBase is the extra black space around the content
	MonitorInternalPaddingBase = 3
)

// GetMonitorOffset returns the total offset for the screen content (including internal padding), scaled
func GetMonitorOffset(scale int) int {
	return (MonitorBezelBase + MonitorPaddingBase + MonitorInternalPaddingBase) * scale
}

// MonitorColor is the classic beige color of an Apple II monitor
var MonitorColor = color.RGBA{0xD8, 0xD0, 0xC0, 0xFF}

// BezelShadowColor is a darker shade for depth
var BezelShadowColor = color.RGBA{0xB8, 0xB0, 0xA0, 0xFF}

// DrawMonitorFrame draws the Apple II monitor frame around the screen
func DrawMonitorFrame(screen *ebiten.Image, width, height int, scale int) {
	// Full monitor background
	vector.DrawFilledRect(screen, 0, 0, float32(width), float32(height), MonitorColor, true)

	bezel := float32(MonitorBezelBase * scale)
	bezelPlusPadding := float32((MonitorBezelBase + MonitorPaddingBase) * scale)

	// Draw an inner "shadow" or depth to the bezel
	vector.StrokeRect(screen, bezel/2, bezel/2, float32(width)-bezel, float32(height)-bezel, float32(2*scale), BezelShadowColor, true)

	// Inner screen area (the "glass" or CRT area)
	// This area is slightly larger than the content to provide the internal padding
	screenAreaWidth := float32(width) - 2*bezelPlusPadding
	screenAreaHeight := float32(height) - 2*bezelPlusPadding
	vector.DrawFilledRect(screen, bezelPlusPadding, bezelPlusPadding, screenAreaWidth, screenAreaHeight, color.Black, true)

	// Power LED (bottom right)
	ledSize := float32(4 * scale)
	// 3 pixels (scaled) from the bezel shadow line (which is at bezel/2)
	// Ensuring the LED is at the same distance from the right and bottom edges
	ledX := float32(width) - bezel/2 - float32(3*scale) - ledSize
	ledY := float32(height) - bezel/2 - float32(3*scale) - ledSize

	vector.DrawFilledRect(screen, ledX, ledY, ledSize, ledSize, color.RGBA{0x00, 0x80, 0x00, 0xFF}, true)         // Dim green
	vector.DrawFilledRect(screen, ledX+ledSize/4, ledY+ledSize/4, ledSize/2, ledSize/2, color.RGBA{0x00, 0xFF, 0x00, 0xFF}, true) // Bright green center

	// Apple II Logo (bottom left)
	// Symmetrical to the LED
	logoImg := getLogoImage()
	if logoImg != nil {
		logoSize := float32(8 * scale)
		logoX := bezel/2 + float32(3*scale)
		logoY := float32(height) - bezel/2 - float32(3*scale) - logoSize

		op := &ebiten.DrawImageOptions{}
		// Scale to fit the target logo size
		sw, sh := logoImg.Bounds().Dx(), logoImg.Bounds().Dy()
		op.GeoM.Scale(float64(logoSize)/float64(sw), float64(logoSize)/float64(sh))
		op.GeoM.Translate(float64(logoX), float64(logoY))
		op.Filter = ebiten.FilterLinear // Smoother scaling for the logo
		screen.DrawImage(logoImg, op)
	}
}
