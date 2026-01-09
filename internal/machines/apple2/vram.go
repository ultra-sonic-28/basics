package apple2

const (
	TextBase = 0x400
	TextCols = 40
	TextRows = 24
	TextSize = 0x400 // 1KB réservé
)

type VRAM struct {
	mem [TextSize]byte
}

func NewVRAM() *VRAM {
	v := &VRAM{}
	v.Clear()
	return v
}

func (v *VRAM) Write(addr int, value byte) {
	offset := addr - TextBase
	if offset < 0 || offset >= len(v.mem) {
		return
	}
	v.mem[offset] = value
}

func (v *VRAM) Read(addr int) byte {
	offset := addr - TextBase
	if offset < 0 || offset >= len(v.mem) {
		return 0
	}
	return v.mem[offset]
}

func (v *VRAM) Clear() {
	for i := range v.mem {
		v.mem[i] = ' ' | 0x80 // Apple II : bit 7 = inverse/normal
	}
}
