package cmd

import (
	"path/filepath"

	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/afero"
)

func resolveDir(dirFlag string, mc *config.MetanaConfig) string {
	if dirFlag != "" {
		return dirFlag
	}
	if mc != nil && mc.Dir != "" {
		return mc.Dir
	}
	return "migrations"
}

func resolveStore(storeFlag string, mc *config.MetanaConfig) string {
	if storeFlag != "" {
		return storeFlag
	}
	if mc != nil && mc.StoreConn != "" {
		return mc.StoreConn
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
