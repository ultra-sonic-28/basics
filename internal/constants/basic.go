// internal/constants/basic.go
package constants

// BASIC types
const (
	BASIC_APPLE byte = 0
	BASIC_C64   byte = 1
	BASIC_AMS   byte = 2
)

// BASIC versions (1 octet chacun)
var BasicVersion = map[byte]byte{
	BASIC_APPLE: 10, // version 1.0
	BASIC_C64:   10, // version 1.0
	BASIC_AMS:   10, // version 1.0
}

// Mapping type → display name
var BasicName = map[byte]string{
	BASIC_APPLE: "APPLE",
	BASIC_C64:   "C64",
	BASIC_AMS:   "AMS6128",
}
