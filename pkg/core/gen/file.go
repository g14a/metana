package gen

import (
	"bytes"
	"fmt"
	"github.com/g14a/metana/pkg"
	"go/format"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/afero"

	tpl2 "github.com/g14a/metana/pkg/core/tpl"

	"github.com/iancoleman/strcase"
)

func CreateMigrationFile(migrationsDir, file string, FS afero.Fs) (string, error) {
	nm := tpl2.NewMigration{
		MigrationName: strcase.ToCamel(file),
		Timestamp:     strconv.Itoa(int(time.Now().Unix())),
	}

	fileName := fmt.Sprintf(migrationsDir+"/scripts/%s-%s.go", nm.Timestamp, nm.MigrationName)

	mainFile, err := FS.Create(fileName)
	if err != nil {
		return "", err
	}

	defer func(mainFile afero.File) {
		err := mainFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(mainFile)

	mainTemplate := template.Must(
		template.New("root").
			Parse(string(tpl2.MigrationTemplate())))

	buff := new(bytes.Buffer)
	err = mainTemplate.Execute(buff, nm)
	if err != nil {
		return "", err
	}

	fmtBytes, err := format.Source(buff.Bytes())
	err = afero.WriteFile(FS, fileName, fmtBytes, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func CreateInitConfig(migrationsDir, goModPath string, FS afero.Fs) error {

	migrationRunFile, err := FS.Create(migrationsDir + "/main.go")
	if err != nil {
		return err
	}

	defer func(migrationRunFile afero.File) {
		err := migrationRunFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(migrationRunFile)

	migrationRunTemplate := template.Must(
		template.New("main").
			Parse(string(tpl2.InitMigrationRunTemplate())))

	params := map[string]interface{}{
		"pwd": goModPath,
		"dir": migrationsDir,
	}

	buff := new(bytes.Buffer)
	err = migrationRunTemplate.Execute(buff, params)

	if err != nil {
		return err
	}

	fmtBytes, err := format.Source(buff.Bytes())
	err = afero.WriteFile(FS, migrationsDir+"/main.go", fmtBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func MigrationExists(migrationsDir, migrationName string, FS afero.Fs) (bool, error) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	migrations, err := pkg.GetMigrations(migrationsDir, FS)
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
