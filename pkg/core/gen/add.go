package gen

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	tpl2 "github.com/g14a/metana/pkg/core/tpl"

	"github.com/g14a/metana/pkg"
	"github.com/iancoleman/strcase"
)

func Regen(migrationsDir, migrationName, fileName string, firstMigration bool) error {
	lower := strcase.ToLowerCamel(migrationName)
	input, err := os.ReadFile(migrationsDir + "/main.go")
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

	err = os.WriteFile(migrationsDir+"/main.go", []byte(output), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("gofmt", "-w", migrationsDir+"/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}

func MigrationExists(migrationsDir, migrationName string) (bool, error) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	migrations, err := pkg.GetMigrations(migrationsDir)
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
