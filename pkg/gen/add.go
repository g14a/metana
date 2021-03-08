package gen

import (
	"fmt"
	"go-migrate/pkg"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
)

func AddMigration(migrationName, fileName string) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	regenerateMain(camelCaseMigration, fileName)
}

func regenerateMain(migrationName, fileName string) {
	lower := strcase.ToLowerCamel(migrationName)

	input, err := ioutil.ReadFile("migrations/main.go")

	lines := strings.Split(string(input), "\n")

	var firstReturn bool

	for i, line := range lines {
		if !firstReturn && strings.Contains(line, "return nil") {
			lines[i] = lower + "Migration := &" + migrationName + "Migration{}\n err" + migrationName + " := " +
				lower + "Migration.Up()\n if err" + migrationName + " != nil {\n return fmt.Errorf(\"" + fileName +
				", %w\", err" + migrationName + ")}\n\n return nil"

			firstReturn = true
		} else if strings.Contains(line, "return nil") {
			lines[i] = lower + "Migration := &" + migrationName + "Migration{}\n err" + migrationName + " := " +
				lower + "Migration.Down()\n if err" + migrationName + " != nil {\n return fmt.Errorf(\"" + fileName +
				", %w\", err" + migrationName + ")}\n\n return nil"
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
