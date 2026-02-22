package common

import (
	"math"
	"strconv"
	"strings"
)

func FormatNumber(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return "OVER FLOW ERROR"
	}

	// Applesoft standard formatting
	absF := math.Abs(f)

	// Applesoft uses scientific notation for abs(f) >= 1e9
	// and for some very small numbers.
	if absF >= 1e9 || (absF > 0 && absF < 0.01) {
		s := strconv.FormatFloat(f, 'e', -1, 64)
		s = strings.ReplaceAll(s, "E", "e")
		// FormatFloat 'e' always puts the sign after 'e'. 
		// If it's positive, it puts '+'. If negative, '-'.
		// Example: 1e+10, 1e-05
		return s
	}

	// Entiers
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}

	// Autres nombres (format standard)
	s := strconv.FormatFloat(f, 'g', -1, 64)
	s = strings.ReplaceAll(s, "E", "e")
	if strings.Contains(s, "e") && !strings.Contains(s, "e-") && !strings.Contains(s, "e+") {
		s = strings.Replace(s, "e", "e+", 1)
	}

	return s
}
