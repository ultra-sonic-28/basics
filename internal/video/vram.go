package video

// VRAM représente une mémoire vidéo virtuelle.
type VRAM interface {
	Write(addr int, value byte)
	Read(addr int) byte
	Clear()
}
