package binary

const MagicString = "BASC"

type Header struct {
	Magic     [4]byte // "BASC"
	BasicType byte
	Version   byte
	NodeCount uint32
	CRC32     uint32
}
