package gen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"

	tpl2 "github.com/g14a/metana/pkg/core/tpl"

	"github.com/iancoleman/strcase"
)

func CreateMigrationFile(migrationsDir, file string) (string, error) {
	nm := tpl2.NewMigration{
		MigrationName: strcase.ToCamel(file),
		Timestamp:     strconv.Itoa(int(time.Now().Unix())),
	}

	fileName := fmt.Sprintf(migrationsDir+"/scripts/%s-%s.go", nm.Timestamp, nm.MigrationName)

	mainFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer mainFile.Close()

	mainTemplate := template.Must(template.New("root").Parse(string(tpl2.MigrationTemplate())))
	err = mainTemplate.Execute(mainFile, nm)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("gofmt", "-w", fileName)
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return fileName, nil
}

func CreateInitConfig(migrationsDir, pwd string) error {

	migrationRunFile, err := os.Create(migrationsDir + "/main.go")
	if err != nil {
		return err
	}

	defer func(migrationRunFile *os.File) {
		err := migrationRunFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(migrationRunFile)

	migrationRunTemplate := template.Must(template.New("main").Parse(string(tpl2.InitMigrationRunTemplate())))
	err = migrationRunTemplate.Execute(migrationRunFile, map[string]interface{}{
		"pwd": pwd,
		"dir": migrationsDir,
	})

	if err != nil {
		return err
	}

	cmd := exec.Command("gofmt", "-w", migrationsDir+"/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}
