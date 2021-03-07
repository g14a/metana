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

	fmt.Println(migrateUpStart, migrateDownStart)
	
	lines[migrateUpEnd] = lower + "Migration := &" + migrationName + "Migration{}\n err" + migrationName + " := " +
		lower + "Migration.Up()\n if err" + migrationName + " != nil {\n return err" + migrationName + "}\n}"

	//lines[migrateUpEnd+migrateDownStart-migrateUpStart+4] = lower + "Migration := &" + migrationName + "Migration{}\n err" + migrationName + " := " +
	//	lower + "Migration.Down()\n if err" + migrationName + " != nil {\n return err" + migrationName + "}\n}"

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
