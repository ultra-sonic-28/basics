package font

import (
	"basics/internal/constants"
	"basics/internal/logger"
)

func DefaultFontForMode(basicType byte) *BitmapFont {
	switch basicType {
	case constants.BASIC_APPLE:
		logger.Info("Use Font7x8")
		return Font7x8
	default:
		logger.Info("Use Font8x8")
		return Font8x8
	}
}
