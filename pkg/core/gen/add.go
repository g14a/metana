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

func Regen(opts RegenOpts) error {
	lower := strcase.ToLowerCamel(opts.MigrationName)
	var input []byte
	if opts.Environment == "" {
		inputBytes, err := afero.ReadFile(opts.FS, opts.MigrationsDir+"/main.go")
		if err != nil {
			return err
		}
		input = inputBytes
	} else {
		inputBytes, err := afero.ReadFile(opts.FS, opts.MigrationsDir+"/environments/"+opts.Environment+"/main.go")
		if err != nil {
			return err
		}
		input = inputBytes
	}

	lines := strings.Split(string(input), "\n")

	var firstReturn bool
	timeStamp := strings.TrimLeft(strings.Split(opts.Filename, "-")[0], "scripts/")

	addMigrationTemplate := template.New("add")

	nm := tpl2.NewMigration{
		Lower:         lower,
		MigrationName: opts.MigrationName,
		Timestamp:     timeStamp,
		Filename:      opts.Filename,
	}

	for i, line := range lines {
		if !firstReturn && strings.Contains(line, "return nil") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl2.AddMigrationTemplate(true)))
			if errAdd != nil {
				return errAdd
			}
			err := addMigrationTemplate.Execute(&tplBuffer, nm)
			if err != nil {
				return err
			}

			lines[i] = tplBuffer.String()
			firstReturn = true
		} else if strings.Contains(line, "func MigrateDown") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl2.AddMigrationTemplate(false)))
			if errAdd != nil {
				return errAdd
			}
			err := addMigrationTemplate.Execute(&tplBuffer, nm)
			if err != nil {
				return err
			}
			if opts.FirstMigration {
				tplBuffer.WriteString("\nreturn nil")
			}
			lines[i+1] = tplBuffer.String()
		}
	}

	output := strings.Join(lines, "\n")

	fmtOutput, err := format.Source([]byte(output))
	if err != nil {
		return err
	}

	if opts.Environment == "" {
		err = afero.WriteFile(opts.FS, opts.MigrationsDir+"/main.go", fmtOutput, 0644)
		if err != nil {
			return err
		}
	} else {
		err = afero.WriteFile(opts.FS, opts.MigrationsDir+"/environments/"+opts.Environment+"/main.go", fmtOutput, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

type RegenOpts struct {
	MigrationsDir  string
	MigrationName  string
	Filename       string
	FirstMigration bool
	Environment    string
	FS             afero.Fs
}
