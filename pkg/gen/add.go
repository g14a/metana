package gen

import (
	"bytes"
	"fmt"
	"github.com/g14a/go-migrate/pkg"
	"github.com/g14a/go-migrate/pkg/tpl"
	"os"
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
	input, err := os.ReadFile("migrations/main.go")
	if err != nil {
		return err
	}
	lines := strings.Split(string(input), "\n")

	var firstReturn bool
	timeStamp := strings.TrimLeft(strings.Split(fileName, "-")[0], "scripts/")

	addMigrationTemplate := template.New("add")

	nm := tpl.NewMigration{
		Lower:         lower,
		MigrationName: migrationName,
		Timestamp:     timeStamp,
		Filename:      fileName,
	}

	for i, line := range lines {
		if !firstReturn && strings.Contains(line, "return nil") {
			var tplBuffer bytes.Buffer
			addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl.AddMigrationTemplate(true)))
			if errAdd != nil {
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
			addMigrationTemplate, errAdd := addMigrationTemplate.Parse(string(tpl.AddMigrationTemplate(false)))
			if errAdd != nil {
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

	err = os.WriteFile("migrations/main.go", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("gofmt", "-w", "migrations/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}

func MigrationExists(migrationName string) (bool, error) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	migrations, err := pkg.GetMigrations()
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
