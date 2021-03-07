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

	var firstReturn bool

	for i, line := range lines {
		if !firstReturn && strings.Contains(line, "return nil") {
			lines[i] = lower + "Migration := &" + migrationName + "Migration{}\n err" + migrationName + " := " +
				lower + "Migration.Up()\n if err" + migrationName + " != nil {\n return err" + migrationName + "}\n\n return nil"

			firstReturn = true
		} else if strings.Contains(line, "return nil") {
			lines[i] = lower + "Migration := &" + migrationName + "Migration{}\n err" + migrationName + " := " +
				lower + "Migration.Down()\n if err" + migrationName + " != nil {\n return err" + migrationName + "}\n\n return nil"
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
