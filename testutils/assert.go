package testutils

import (
	"fmt"
	"testing"
)

// Equal échoue le test si got != want
func Equal[T comparable](t *testing.T, msg string, got T, want T) {
	t.Helper()

	if got != want {
		if len(msg) == 0 {
			msg = "assert.Equal failed:"
		}

		// Convertir got et want en string pour %q
		gotStr := fmt.Sprintf("%v", got)
		wantStr := fmt.Sprintf("%v", want)

		t.Fatalf("%s got=%q want=%q", msg, gotStr, wantStr)
	}

	RecordAssertion(t)
}

// NotEqual échoue le test si got == want
func NotEqual[T comparable](t *testing.T, msg string, got T, want T) {
	t.Helper()

	if got == want {
		if len(msg) == 0 {
			msg = "assert.NotEqual failed:"
		}

		// Convertir got et want en string pour %q
		gotStr := fmt.Sprintf("%v", got)

		t.Fatalf("%s value=%q", msg, gotStr)
	}

	RecordAssertion(t)
}

// True échoue si value != true
func True(t *testing.T, msg string, value bool) {
	t.Helper()

	if !value {
		if len(msg) == 0 {
			msg = "assert.True failed:"
		}
		t.Fatalf("%s", msg)
	}

	RecordAssertion(t)
}

// False échoue si value != false
func False(t *testing.T, msg string, value bool) {
	t.Helper()

	if value {
		if len(msg) == 0 {
			msg = "assert.False failed:"
		}
		t.Fatalf("%s", msg)
	}

	RecordAssertion(t)
}

// AssertMapContainsKeys échoue si une ou plusieurs clés sont absentes du map.
//
// Exemple :
//
//	AssertMapContainsKeys(t, "missing keywords", myMap, []string{"A", "B"})
func AssertMapContainsKeys[K comparable, V any](
	t *testing.T,
	msg string,
	m map[K]V,
	expected []K,
) {
	t.Helper()

	for _, key := range expected {
		if _, ok := m[key]; !ok {
			if msg == "" {
				msg = "assert.MapContainsKeys failed:"
			}
			t.Fatalf("%s missing key=%v", msg, key)
		}
	}

	RecordAssertion(t)
}
