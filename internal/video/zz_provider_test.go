package video

import (
	"testing"

	"basics/testutils"
)

// MockProvider implémente l'interface Provider pour tests
type MockProvider struct {
	modes       []ModeID
	modeInfos   map[ModeID]ModeInfo
	defaultMode ModeID
}

func (m *MockProvider) Modes() []ModeID {
	return m.modes
}

func (m *MockProvider) ModeInfo(id ModeID) (ModeInfo, bool) {
	info, ok := m.modeInfos[id]
	return info, ok
}

func (m *MockProvider) DefaultMode() ModeID {
	return m.defaultMode
}

func TestMockProvider(t *testing.T) {
	mock := &MockProvider{
		modes: []ModeID{
			ModeID("text40"),
			ModeID("text80"),
			ModeID("gfx320"),
		},
		modeInfos: map[ModeID]ModeInfo{
			ModeID("text40"): {
				ID:     "text40",
				Name:   "Text 40x24",
				Width:  40,
				Height: 24,
				Text:   true,
			},
			ModeID("gfx320"): {
				ID:     "gfx320",
				Name:   "Graphics 320x200",
				Width:  320,
				Height: 200,
				Text:   false,
			},
		},
		defaultMode: ModeID("text40"),
	}

	// --- Test Modes ---
	gotModes := mock.Modes()
	testutils.DeepEqual(t, "Modes", gotModes, []ModeID{
		"text40",
		"text80",
		"gfx320",
	})

	// --- Test ModeInfo existant ---
	info, ok := mock.ModeInfo(ModeID("text40"))
	testutils.True(t, "ModeInfo found", ok)
	testutils.Equal(t, "ModeInfo.ID", info.ID, "text40")
	testutils.Equal(t, "ModeInfo.Name", info.Name, "Text 40x24")
	testutils.Equal(t, "ModeInfo.Width", info.Width, 40)
	testutils.Equal(t, "ModeInfo.Height", info.Height, 24)
	testutils.True(t, "ModeInfo.Text", info.Text)

	// --- Test ModeInfo inexistant ---
	_, ok = mock.ModeInfo(ModeID("nonexistent"))
	testutils.False(t, "ModeInfo not found", ok)

	// --- Test DefaultMode ---
	def := mock.DefaultMode()
	testutils.Equal(t, "DefaultMode", def, ModeID("text40"))
}
