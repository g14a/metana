package gen

import (
	"bytes"
	"go/format"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	tpl2 "github.com/g14a/metana/pkg/core/tpl"

	"github.com/iancoleman/strcase"
)

func Regen(migrationsDir, migrationName, fileName string, firstMigration bool, FS afero.Fs) error {
	lower := strcase.ToLowerCamel(migrationName)
	input, err := afero.ReadFile(FS, migrationsDir+"/main.go")
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), "\n")

	var firstReturn bool
	timeStamp := strings.TrimLeft(strings.Split(fileName, "-")[0], "scripts/")

	addMigrationTemplate := template.New("add")

	nm := tpl2.NewMigration{
		Lower:         lower,
		MigrationName: migrationName,
		Timestamp:     timeStamp,
		Filename:      fileName,
	}

	for i, line := range lines {
		if !firstReturn && strings.Contains(line, "return nil") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl2.AddMigrationTemplate(true)))
			if errAdd != nil {
				return err
			}
			err = addMigrationTemplate.Execute(&tplBuffer, nm)
			if err != nil {
				return err
			}

			lines[i] = tplBuffer.String()
			firstReturn = true
		} else if strings.Contains(line, "func MigrateDown") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl2.AddMigrationTemplate(false)))
			if errAdd != nil {
				return err
			}
			err = addMigrationTemplate.Execute(&tplBuffer, nm)
			if err != nil {
				return err
			}
			if firstMigration {
				tplBuffer.WriteString("\nreturn nil")
			}
			lines[i+1] = tplBuffer.String()
		}
	}

	output := strings.Join(lines, "\n")

	fmtOutput, err := format.Source([]byte(output))

	err = afero.WriteFile(FS, migrationsDir+"/main.go", fmtOutput, 0644)
	if err != nil {
		return err
	}

	return nil
}
