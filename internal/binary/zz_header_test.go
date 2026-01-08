package binary_test

import (
	"basics/internal/binary"
	"basics/testutils"
	"testing"
)

func TestHeaderMagic(t *testing.T) {
	// Vérifier la constante MagicString
	testutils.Equal(t, "MagicString", binary.MagicString, "BASC")

	// Création d'un header vide
	h := binary.Header{}

	// Vérifier que Magic est bien initialisé vide
	var emptyMagic [4]byte
	testutils.Equal(t, "Magic initial empty", h.Magic, emptyMagic)

	// On peut remplir Magic et vérifier
	copy(h.Magic[:], []byte(binary.MagicString))
	testutils.Equal(t, "Magic filled", string(h.Magic[:]), binary.MagicString)

	// Vérifier les autres champs par défaut
	testutils.Equal(t, "BasicType default", h.BasicType, byte(0))
	testutils.Equal(t, "Version default", h.Version, byte(0))
	testutils.Equal(t, "NodeCount default", h.NodeCount, uint32(0))
	testutils.Equal(t, "CRC32 default", h.CRC32, uint32(0))
}
