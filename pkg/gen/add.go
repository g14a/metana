package gen

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
)

func AddMigration(migrationName string) {
	camelCaseMigration := strcase.ToCamel(migrationName)

	regenerateMain(camelCaseMigration)
}

func regenerateMain(migrationName string) {
	lower := strcase.ToLowerCamel(migrationName)

	input, err := ioutil.ReadFile("migrations/main.go")

	lines := strings.Split(string(input), "\n")

	var migrateUpStart, migrateDownStart, migrateUpEnd int
	var firstCloseBrace bool

	for i, line := range lines {
		if strings.Contains(line, "func MigrateUp() {") {
			migrateUpStart = i
		}

		if strings.Contains(line, "func MigrateDown() {") {
			migrateDownStart = i
		}

		if !firstCloseBrace && !strings.Contains(line, "Migration{") && strings.Contains(line, "}") {
			firstCloseBrace = true
			migrateUpEnd = i
		}
	}

	lines[migrateUpEnd] = "\nvar " + lower + "Migration _interface.Migration = &" + migrationName + "Migration{}\n" + lower + "Migration.Up()\n}"

	lines[migrateUpEnd+migrateDownStart-migrateUpStart] = "\nvar " + lower + "Migration _interface.Migration = &" + migrationName + "Migration{}\n" + lower + "Migration.Down()\n}"

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
