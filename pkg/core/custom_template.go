package core

import (
	"strings"

	"github.com/spf13/afero"
)

func ParseCustomTemplate(wd string, fileName string, FS afero.Fs) (string, string) {
	bytes, err := afero.ReadFile(FS, wd+"/"+fileName)
	if err != nil {
		return "", ""
	}

	lines := strings.Split(string(bytes), "\n")

	var upBuilder, downBuilder strings.Builder
	for i, line := range lines {
		if strings.Contains(line, "Up() error") {
			for k := i + 1; k < len(lines); k++ {
				if strings.Contains(lines[k], "}") {
					break
				}
				upBuilder.WriteString(lines[k] + "\n")
			}
		}
		if strings.Contains(line, "Down() error") {
			for k := i + 1; k < len(lines); k++ {
				if strings.Contains(lines[k], "}") {
					break
				}
				downBuilder.WriteString(lines[k] + "\n")
			}
		}
	}

	return strings.TrimSpace(upBuilder.String()), strings.TrimSpace(downBuilder.String())
}
