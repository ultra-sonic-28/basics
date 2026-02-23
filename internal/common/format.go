package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func FormatNumber(f float64) string {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return "OVER FLOW ERROR"
	}

	if f == 0 {
		return "0"
	}

	absF := math.Abs(f)

	// Applesoft uses scientific notation for abs(f) >= 1e9
	// and for very small numbers. User says from 3e-7 it should be scientific.
	// 3e-7 is 0.0000003. So threshold 1e-6 (0.000001) would make 3e-7 scientific.
	if absF >= 1e9 || absF < 0.000001 {
		s := strconv.FormatFloat(f, 'e', -1, 64)
		s = strings.ReplaceAll(s, "E", "e")
		return s
	}

	// Entiers
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}

	// Autres nombres (format standard)
	// On utilise 'f' pour forcer le format décimal
	s := strconv.FormatFloat(f, 'f', -1, 64)

	return s
}

// ParseFloatApplesoft parses a string into a float64 according to Applesoft VAL() rules.
func ParseFloatApplesoft(s string) (float64, error) {
	// Skip leading spaces
	s = strings.TrimLeft(s, " ")
	if s == "" {
		return 0, nil
	}

	// First character must be a digit, +, -, or .
	if !unicode.IsDigit(rune(s[0])) && s[0] != '.' && s[0] != '+' && s[0] != '-' {
		return 0, nil
	}

	// Find the end of the numeric part
	// It stops at the first character that cannot be part of the number.
	// Allowed: digits, one dot, 'e' or 'E' (scientific notation), and + or - after 'e'.
	end := 0
	hasDot := false
	hasExponent := false
	
	for i, r := range s {
		if unicode.IsDigit(r) {
			end = i + 1
			continue
		}
		if r == '.' && !hasDot && !hasExponent {
			hasDot = true
			end = i + 1
			continue
		}
		if (r == 'e' || r == 'E') && !hasExponent {
			hasExponent = true
			end = i + 1
			continue
		}
		if (r == '+' || r == '-') {
			if i == 0 {
				end = i + 1
				continue
			}
			// Allowed after 'e'
			if i > 0 && (s[i-1] == 'e' || s[i-1] == 'E') {
				end = i + 1
				continue
			}
		}
		// Any other character stops the parsing
		break
	}

	if end == 0 {
		return 0, nil
	}

	numStr := s[:end]
	
	// If it's just "+" or "-" or "." or "e" it should be 0
	checkStr := strings.ToLower(numStr)
	if checkStr == "+" || checkStr == "-" || checkStr == "." || checkStr == "e" || checkStr == "e+" || checkStr == "e-" {
		return 0, nil
	}

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		// Check for overflow
		if strings.Contains(err.Error(), "value out of range") {
			return 0, fmt.Errorf("OVER FLOW ERROR")
		}
		return 0, nil
	}

	if math.IsInf(val, 0) {
		return 0, fmt.Errorf("OVER FLOW ERROR")
	}

	return val, nil
}
