package gen

import (
	"bytes"
	"github.com/g14a/metana/pkg"
	"go/format"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/afero"

	tpl2 "github.com/g14a/metana/pkg/core/tpl"

	"github.com/iancoleman/strcase"
)

func Regen(opts RegenOpts) error {
	lower := strcase.ToLowerCamel(opts.MigrationName)

	addMigrationTemplate := template.New("add")

	var migrationsToCreate []tpl2.NewMigration

	for _, m := range opts.Migrations {
		migrationName := strings.TrimLeftFunc(m.Name, func(r rune) bool {
			return unicode.IsNumber(r) || r == '-'
		})
		ts := strings.TrimRightFunc(m.Name, func(r rune) bool {
			return !unicode.IsNumber(r)
		})
		migrationsToCreate = append(migrationsToCreate, tpl2.NewMigration{
			Lower:         strcase.ToLowerCamel(strings.TrimRight(migrationName, ".go")),
			MigrationName: strings.TrimRight(migrationName, ".go"),
			Timestamp:     ts,
			Filename:      m.Name,
		})
	}

	timeStamp := strings.TrimLeft(strings.Split(opts.Filename, "-")[0], "scripts/")

	migrationsToCreate = append(migrationsToCreate, tpl2.NewMigration{
		Lower:         lower,
		MigrationName: opts.MigrationName,
		Timestamp:     timeStamp,
		Filename:      opts.Filename,
	})

	addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl2.NewAddMigrationTemplate()))
	if errAdd != nil {
		return errAdd
	}

	var tplBuffer bytes.Buffer

	params := Params{
		Pwd:    opts.GoModPath,
		Dir:    opts.MigrationsDir,
		Create: migrationsToCreate,
	}

	err := addMigrationTemplate.Execute(&tplBuffer, params)
	if err != nil {
		return err
	}

	fmtOutput, err := format.Source(tplBuffer.Bytes())
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
	GoModPath      string
	Migrations     []pkg.Migration
	FS             afero.Fs
}

type Params struct {
	Pwd    string
	Dir    string
	Create []tpl2.NewMigration
}
