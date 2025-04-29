package cmd

import (
	"path/filepath"

	"github.com/spf13/afero"
)

// resolveDir simply returns the CLI flag if set, otherwise default "migrations"
func resolveDir() string {
	return "migrations"
}

// resolveStore simply returns the CLI store flag if set, otherwise empty
func resolveStore(storeFlag string) string {
	if storeFlag != "" {
		return storeFlag
	}
	return ""
}

// cleanFinalDir normalizes finalDir to be relative to wd
func cleanFinalDir(wd, dir string) string {
	if filepath.IsAbs(dir) {
		if rel, err := filepath.Rel(wd, dir); err == nil {
			return rel
		}
	}
	return dir
}

// mkdirScripts creates the scripts/ folder safely
func mkdirScripts(fs afero.Fs, baseDir string) error {
	return fs.MkdirAll(filepath.Join(baseDir, "scripts"), 0755)
}
