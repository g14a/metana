package gen

import (
	"bytes"
	"fmt"
	"go-migrate/pkg"
	"go-migrate/pkg/tpl"
	"io/ioutil"
	"os/exec"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

func AddMigration(migrationName, fileName string) error {
	camelCaseMigration := strcase.ToCamel(migrationName)

	return regenerateMain(camelCaseMigration, fileName)
}

func regenerateMain(migrationName, fileName string) error {
	lower := strcase.ToLowerCamel(migrationName)
	input, err := ioutil.ReadFile("migrations/main.go")

	lines := strings.Split(string(input), "\n")

	var firstReturn bool
	timeStamp := strings.Split(fileName, "-")

	addMigrationTemplate := template.New("add")
	nm := tpl.NewMigration{
		Lower:         lower,
		MigrationName: migrationName,
		Timestamp:     timeStamp[0],
		Filename:      fileName,
	}

	for i, line := range lines {
		if !firstReturn && strings.Contains(line, "return nil") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, err = addMigrationTemplate.Parse(string(tpl.AddMigrationTemplate(true)))
			if err != nil {
				return err
			}
			err = addMigrationTemplate.Execute(&tplBuffer, nm)
			if err != nil {
				return err
			}

			lines[i] = tplBuffer.String()

			firstReturn = true
		} else if strings.Contains(line, "return nil") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, err = addMigrationTemplate.Parse(string(tpl.AddMigrationTemplate(false)))
			if err != nil {
				return err
			}
			err = addMigrationTemplate.Execute(&tplBuffer, nm)
			if err != nil {
				return err
			}
			lines[i] = tplBuffer.String()
		}
	}

	output := strings.Join(lines, "\n")

	err = ioutil.WriteFile("migrations/main.go", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("gofmt", "-w", "migrations/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}

func MigrationExists(migrationName string) bool {
	camelCaseMigration := strcase.ToCamel(migrationName)

	migrations, err := pkg.GetMigrations()
	if err != nil {
		fmt.Println(err)
	}

	for _, m := range migrations {
		mig := strings.TrimSuffix(m.Name, ".go")
		mig = strings.TrimLeftFunc(mig, func(r rune) bool {
			return r >= 48 && r <= 57 || r == '-'
		})
		if camelCaseMigration == mig {
			return true
		}
	}

	return false
}
