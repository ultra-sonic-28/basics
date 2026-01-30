package font

import (
	"basics/internal/constants"
	"basics/testutils"
	"testing"
)

func TestDefaultFontForMode(t *testing.T) {
	tests := []struct {
		name      string
		basicType byte
		want      *BitmapFont
	}{
		{
			name:      "APPLE BASIC uses Font7x8",
			basicType: constants.BASIC_APPLE,
			want:      Font7x8,
		},
		{
			name:      "Other BASIC uses Font8x8",
			basicType: constants.BASIC_C64,
			want:      Font8x8,
		},
		{
			name:      "Unknown BASIC defaults to Font8x8",
			basicType: 0xFF,
			want:      Font8x8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultFontForMode(tt.basicType)
			testutils.Equal(t, "BitmapFont pointer", got, tt.want)
		})
	}
}
