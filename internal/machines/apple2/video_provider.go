package apple2

import "basics/internal/video"

const (
	Text40 video.VideoModeID = iota
	LoRes
	HiRes
)

var modes = map[video.VideoModeID]video.ModeInfo{
	Text40: {
		Name:   "Text 40",
		Width:  40,
		Height: 24,
		Text:   true,
	},
	LoRes: {
		Name:   "Lo-Res Graphics",
		Width:  40,
		Height: 48,
	},
	HiRes: {
		Name:   "Hi-Res Graphics",
		Width:  280,
		Height: 192,
	},
}

type VideoProvider struct{}

func NewVideoProvider() *VideoProvider {
	return &VideoProvider{}
}

func (VideoProvider) Modes() []video.VideoModeID {
	return []video.VideoModeID{Text40, LoRes, HiRes}
}

func (VideoProvider) ModeInfo(id video.VideoModeID) (video.ModeInfo, bool) {
	info, ok := modes[id]
	return info, ok
}

func (VideoProvider) DefaultMode() video.VideoModeID {
	return Text40
}
