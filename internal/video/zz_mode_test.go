package video

import (
	"basics/testutils"
	"testing"
)

// MockMode implémente l'interface Mode pour les tests
type MockMode struct {
	info        ModeInfo
	resetCalled bool
}

func (m *MockMode) Info() ModeInfo {
	return m.info
}

func (m *MockMode) Reset() {
	m.resetCalled = true
}

func TestModeInfo_Positive(t *testing.T) {
	info := ModeInfo{
		ID:     "text40",
		Name:   "Text 40x24",
		Width:  40,
		Height: 24,
		Text:   true,
	}

	testutils.Equal(t, "ID", info.ID, "text40")
	testutils.Equal(t, "Name", info.Name, "Text 40x24")
	testutils.Equal(t, "Width", info.Width, 40)
	testutils.Equal(t, "Height", info.Height, 24)
	testutils.True(t, "Text", info.Text)
}

func TestMockMode_Interface(t *testing.T) {
	info := ModeInfo{
		ID:     "gfx320",
		Name:   "Graphics 320x200",
		Width:  320,
		Height: 200,
		Text:   false,
	}

	mode := &MockMode{
		info: info,
	}

	// --- Test Info() ---
	gotInfo := mode.Info()
	testutils.DeepEqual(t, "Info()", gotInfo, info)

	// --- Test Reset() ---
	testutils.False(t, "resetCalled before Reset", mode.resetCalled)
	mode.Reset()
	testutils.True(t, "resetCalled after Reset", mode.resetCalled)
}

func TestModeCollection(t *testing.T) {
	// On peut simuler plusieurs modes
	modes := []Mode{
		&MockMode{
			info: ModeInfo{
				ID:     "text40",
				Name:   "Text 40x24",
				Width:  40,
				Height: 24,
				Text:   true,
			},
		},
		&MockMode{
			info: ModeInfo{
				ID:     "gfx320",
				Name:   "Graphics 320x200",
				Width:  320,
				Height: 200,
				Text:   false,
			},
		},
	}

	testutils.Equal(t, "number of modes", len(modes), 2)

	// Vérifier les infos de chaque mode
	testutils.Equal(t, "first mode ID", modes[0].Info().ID, "text40")
	testutils.Equal(t, "second mode ID", modes[1].Info().ID, "gfx320")
}
