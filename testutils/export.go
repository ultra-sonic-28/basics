package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AssertStats struct {
	PerPackage map[string]int `json:"per_package"`
	Total      int            `json:"total"`
}

func ExportIfRequested() {
	dir := os.Getenv("ASSERT_STATS_DIR")
	if dir == "" {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	pkg := filepath.Base(wd)

	assertCountsMu.Lock()
	defer assertCountsMu.Unlock()

	total := 0
	for _, n := range assertCounts {
		total += n
	}

	stats := AssertStats{
		PerPackage: assertCounts,
		Total:      total,
	}

	filename := filepath.Join(dir, "asserts."+pkg+".json")
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	_ = enc.Encode(stats)
}
