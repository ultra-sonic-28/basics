package runtime

// Video est l’API vidéo exposée à l’interpréteur BASIC
type Video interface {
	Home()
	Print(string)
	HTab(int)
	VTab(int)
}
