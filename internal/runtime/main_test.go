package runtime

import (
	"testing"

	"basics/testutils"
)

func TestMain(m *testing.M) {
	testutils.RunWithAssertTracking(m)
}
