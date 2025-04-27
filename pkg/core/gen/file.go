package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/core/tpl"
	"github.com/iancoleman/strcase"
	"github.com/spf13/afero"
)

func CreateMigrationFile(opts CreateMigrationOpts) (string, error) {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	migrationName := strcase.ToCamel(opts.File)
	fileName := fmt.Sprintf("%s/scripts/%s_%s.go", opts.MigrationsDir, timestamp, opts.File)

	// Use standalone template
	mainTemplate := template.Must(
		template.New("standalone").
			Parse(string(tpl.StandaloneMigrationTemplate())),
	)

	templateData := map[string]string{
		"MigrationName": migrationName,
		"Timestamp":     timestamp,
		"Filename":      fmt.Sprintf("%s_%s.go", timestamp, opts.File),
	}

	var buff bytes.Buffer
	if err := mainTemplate.Execute(&buff, templateData); err != nil {
		return "", err
	}

	fmtBytes, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	if err := afero.WriteFile(opts.FS, fileName, fmtBytes, 0644); err != nil {
		return "", err
	}

	return fileName, nil
}

func MigrationExists(wd, migrationsDir, migrationName string, FS afero.Fs) (bool, error) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	migrations, err := pkg.GetMigrations(wd, migrationsDir, FS)
	if err != nil {
		return false, err
	}

	for _, m := range migrations {
		mig := strings.TrimSuffix(m.Name, ".go")
		mig = strings.TrimLeftFunc(mig, func(r rune) bool {
			return r >= 48 && r <= 57 || r == '-'
		})
		if camelCaseMigration == mig {
			return true, nil
		}
	}

	return false, nil
}

type CreateMigrationOpts struct {
	Wd            string
	MigrationsDir string
	File          string
	FS            afero.Fs
}
