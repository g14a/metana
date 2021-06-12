package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/g14a/metana/pkg/core"

	"github.com/g14a/metana/pkg"

	"github.com/spf13/afero"

	tpl2 "github.com/g14a/metana/pkg/core/tpl"

	"github.com/iancoleman/strcase"
)

func CreateMigrationFile(opts CreateMigrationOpts) (string, error) {
	nm := tpl2.NewMigration{
		MigrationName: strcase.ToCamel(opts.File),
		Timestamp:     strconv.Itoa(int(time.Now().Unix())),
		FirstChar:     string(opts.File[0]),
	}

	var fileName string
	if opts.Environment == "" {
		fileName = fmt.Sprintf(opts.MigrationsDir+"/scripts/%s-%s.go", nm.Timestamp, nm.MigrationName)
	} else {
		fileName = fmt.Sprintf(opts.MigrationsDir+"/environments/"+opts.Environment+"/scripts/%s-%s.go", nm.Timestamp, nm.MigrationName)
	}

	mainFile, err := opts.FS.Create(fileName)
	if err != nil {
		return "", err
	}

	defer func(mainFile afero.File) {
		err := mainFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(mainFile)

	upBuilder, downBuilder := core.ParseCustomTemplate(opts.Wd, opts.CustomTmpl, opts.FS)

	mainTemplate := template.Must(
		template.New("root").
			Parse(string(tpl2.MigrationTemplate(upBuilder, downBuilder))))

	buff := new(bytes.Buffer)
	err = mainTemplate.Execute(buff, nm)
	if err != nil {
		return "", err
	}

	fmtBytes, err := format.Source(buff.Bytes())
	if err != nil {
		return "", err
	}

	err = afero.WriteFile(opts.FS, fileName, fmtBytes, 0644)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func CreateInitConfig(migrationsDir, goModPath string, FS afero.Fs, environment string) error {
	var mrFile afero.File

	if environment == "" {
		migrationRunFile, err := FS.Create(migrationsDir + "/main.go")
		if err != nil {
			return err
		}
		mrFile = migrationRunFile
	} else {
		migrationRunFile, err := FS.Create(migrationsDir + "/environments/" + environment + "/main.go")
		if err != nil {
			return err
		}
		mrFile = migrationRunFile
	}

	defer func(migrationRunFile afero.File) {
		err := migrationRunFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(mrFile)

	migrationRunTemplate := template.Must(
		template.New("main").
			Parse(string(tpl2.InitMigrationRunTemplate())))

	var params map[string]interface{}

	if environment == "" {
		params = map[string]interface{}{
			"pwd": goModPath,
			"dir": migrationsDir,
		}
	} else {
		params = map[string]interface{}{
			"pwd": goModPath,
			"dir": migrationsDir + "/environments/" + environment,
		}
	}

	buff := new(bytes.Buffer)
	err := migrationRunTemplate.Execute(buff, params)
	if err != nil {
		return err
	}

	fmtBytes, err := format.Source(buff.Bytes())
	if err != nil {
		return err
	}
	if environment == "" {
		err = afero.WriteFile(FS, migrationsDir+"/main.go", fmtBytes, 0644)
		if err != nil {
			return err
		}
	} else {
		err = afero.WriteFile(FS, migrationsDir+"/environments/"+environment+"/main.go", fmtBytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func MigrationExists(wd, migrationsDir, migrationName string, FS afero.Fs, environment string) (bool, error) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	migrations, err := pkg.GetMigrations(wd, migrationsDir, FS, environment)
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
	CustomTmpl    string
	Environment   string
	FS            afero.Fs
}
