package environments

import (
	"os"

	"github.com/spf13/afero"
)

func CheckExistingMigrationSetup(FS afero.Fs, wd string) bool {
	_, err := FS.Stat(wd + "/migrations/main.go")
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckExistingEnvironment(FS afero.Fs, wd string, dir, env string) bool {
	r, err := FS.Stat(wd + "/" + dir + "/environments/" + env)
	if os.IsNotExist(err) {
		return false
	}
	if r.IsDir() {
		return true
	}
	return false
}
