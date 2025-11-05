package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var writeMu sync.Mutex

// WriteJSONAtomically writes the given value as pretty JSON to outputDir/filename.
func WriteJSONAtomically(outputDir, filename string, v interface{}) error {
	writeMu.Lock()
	defer writeMu.Unlock()

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	tmp := filepath.Join(outputDir, filename+".tmp")
	out := filepath.Join(outputDir, filename)
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, out)
}
