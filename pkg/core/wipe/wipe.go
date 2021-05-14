package wipe

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/g14a/metana/pkg"
	"github.com/g14a/metana/pkg/initpkg"
	s "github.com/g14a/metana/pkg/store"
	"github.com/iancoleman/strcase"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Wipe(migrationsDir string, storeConn string) error {
	store, err := s.GetStoreViaConn(storeConn, migrationsDir)
	if err != nil {
		return err
	}

	track, err := store.Load()
	if err != nil {
		return err
	}

	if len(track.Migrations) == 0 {
		color.Yellow("No migrations found to wipe.\nTry creating them or running existing ones.")
	}

	for _, m := range track.Migrations {
		fmt.Println(migrationsDir + "/scripts/" + m.Title)
		err := os.Remove(migrationsDir + "/scripts/" + m.Title)
		if err != nil {
			return err
		}
	}

	migrations, err := pkg.GetMigrations(migrationsDir)
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		color.Yellow("No migrations found to wipe.")
		return nil
	}

	var mainBuilder strings.Builder

	mainImportsComponent := getMainAndImportsComponent(migrationsDir)
	mainBuilder.Write(mainImportsComponent)

	mainBuilder.WriteString("\nfunc MigrateUp(upUntil string, lastRunTS int) (err error){\n")

	for _, m := range migrations {
		upComp := getMigrateComponent(m, true)
		mainBuilder.Write(upComp)
	}
	mainBuilder.WriteString("\nreturn nil\n}")

	mainBuilder.WriteString("\n\nfunc MigrateDown(downUntil string, lastRunTS int) (err error){\n")

	for _, m := range migrations {
		upComp := getMigrateComponent(m, false)
		mainBuilder.Write(upComp)
	}
	mainBuilder.WriteString("\nreturn nil\n}\n\n")

	mainmainComponent := getMainOfMain()

	mainBuilder.Write(mainmainComponent)

	err = os.WriteFile(migrationsDir+"/main.go", []byte(mainBuilder.String()), 0644)
	if err != nil {
		fmt.Println(err)
	}

	cmd := exec.Command("gofmt", "-w", migrationsDir+"/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}

func getMainAndImportsComponent(migrationsDir string) []byte {
	goModPath, err := initpkg.GetGoModPath()
	if err != nil {
		fmt.Println(err)
	}

	return []byte(`// This file is auto generated. DO NOT EDIT!
	package main
	
	import (
		"flag"
		"fmt"
		"os"` +
		"\n\"" +
		goModPath + "/" + migrationsDir + "/scripts\"\n" +
		")")
}

func getMigrateComponent(m pkg.Migration, up bool) []byte {
	ts, migrationName, err := pkg.GetComponents(m.Name)
	if err != nil {
		fmt.Println(err)
	}

	lowerMigration := strcase.ToLowerCamel(migrationName)

	if up {
		return []byte("\n" + lowerMigration + "Migration := &scripts." + migrationName + "Migration{}\n" +
			lowerMigration + "Migration.Timestamp = " + strconv.Itoa(ts) + "\n" +
			lowerMigration + "Migration.Filename = \"" + m.Name + "\"\n" +
			lowerMigration + "Migration.MigrationName = \"" + migrationName + "\"\n\n" +
			"if lastRunTS < " + lowerMigration + "Migration.Timestamp {\n" +
			"fmt.Printf(\"  >>> Migrating up: %s" + `\n` + "\", " + lowerMigration + "Migration.Filename)\n" +
			"err" + migrationName + " := " + lowerMigration + "Migration.Up()\n\n" +
			"if err" + migrationName + " != nil {\n" +
			"fmt.Errorf(\"%w\", err" + migrationName + ")\n" +
			"}\n\n" +
			"fmt.Fprintf(os.Stderr, \"" + m.Name + `\n")` +
			"\n}\n\n" +
			"if upUntil == \"" + migrationName + "\" {\n" +
			"if lastRunTS == " + lowerMigration + "Migration.Timestamp {\nreturn\n}\n" +
			"fmt.Printf(\"  >>> Migrated up until: %s" + `\n",` + lowerMigration + "Migration.Filename)\nreturn\n}\n")
	}

	return []byte("\n" + lowerMigration + "Migration := &scripts." + migrationName + "Migration{}\n" +
		lowerMigration + "Migration.Timestamp = " + strconv.Itoa(ts) + "\n" +
		lowerMigration + "Migration.Filename = \"" + m.Name + "\"\n" +
		lowerMigration + "Migration.MigrationName = \"" + migrationName + "\"\n\n" +
		"if lastRunTS >= " + lowerMigration + "Migration.Timestamp {\n" +
		"fmt.Printf(\"  >>> Migrating down: %s" + `\n` + "\", " + lowerMigration + "Migration.Filename)\n" +
		"err" + migrationName + " := " + lowerMigration + "Migration.Down()\n\n" +
		"if err" + migrationName + " != nil {\n" +
		"fmt.Errorf(\"%w\", err" + migrationName + ")\n" +
		"}\n\n" +
		"fmt.Fprintf(os.Stderr, \"" + m.Name + `\n")` +
		"\n}\n\n" +
		"if downUntil == \"" + migrationName + "\" {\n" +
		"if lastRunTS == " + lowerMigration + "Migration.Timestamp {\nreturn\n}\n" +
		"fmt.Printf(\"  >>> Migrated down until: %s" + `\n",` + lowerMigration + "Migration.Filename)\nreturn\n}\n\n")
}

func getMainOfMain() []byte {
	return []byte(`func main() {
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)

	var upUntil, downUntil string
	var lastRunTS int
	upCmd.StringVar(&upUntil, "until", "", "")
	upCmd.IntVar(&lastRunTS, "last-run-ts", 0, "")
	downCmd.StringVar(&downUntil, "until", "", "")
	downCmd.IntVar(&lastRunTS, "last-run-ts", 0, "")

	switch os.Args[1] {
	case "up":
		err := upCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
	case "down":
		err := downCmd.Parse(os.Args[2:])
		if err != nil {
			return
		}
	}

	if upCmd.Parsed() {
		err := MigrateUp(upUntil, lastRunTS)
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
		}
	}

	if downCmd.Parsed() {
		err := MigrateDown(downUntil, lastRunTS)
		if err != nil {
			fmt.Fprintf(os.Stdout, err.Error())
		}
	}
}
`)
}
