package gen

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"strings"
)

func AddMigration(migrationName string)  {
	camelCaseMigration := strcase.ToCamel(migrationName)

	Parse(camelCaseMigration)
}

func Parse(migrationName string) {
	lower := strcase.ToLowerCamel(migrationName)

	input, err := ioutil.ReadFile("migrations/main.go")

	lines := strings.Split(string(input), "\n")

	var migrateUpEnd int
	var firstCloseBrace bool

	for i, line := range lines {
		if !firstCloseBrace && !strings.Contains(line, "Migration{") && strings.Contains(line, "}") {
			firstCloseBrace = true
			migrateUpEnd = i
		}
	}

	lines[migrateUpEnd] = "var " + lower + "Migration _interface.Migration = &" + migrationName + "Migration{}\n\n" + lower + "Migration.Up()\n}"

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("migrations/main.go", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

/*
func MigrateUp() {
	var j _interface.Migration = &ProbationDateMigration{}

	j.Up()
}
*/

/*
package main

import (
	"os"
)

func MigrateUp() {

}

func MigrateDown() {

}

func main() {
	if os.Args[0] == "up" {
		MigrateUp()
	}

	if os.Args[0] == "down" {
		MigrateDown()
	}
}
 */