package gen

import (
	"fmt"
	"github.com/g14a/go-migrate/pkg/tpl"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
)

func CreateMigrationFile(migrationsDir, file string) (string, error) {
	nm := tpl.NewMigration{
		MigrationName: strcase.ToCamel(file),
		Timestamp:     strconv.Itoa(int(time.Now().Unix())),
	}

	fileName := fmt.Sprintf(migrationsDir+"/scripts/%s-%s.go", nm.Timestamp, nm.MigrationName)

	mainFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer mainFile.Close()

	mainTemplate := template.Must(template.New("root").Parse(string(tpl.MigrationTemplate())))
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

	storeFile, err := os.Create(migrationsDir + "/store.go")
	if err != nil {
		return err
	}

	defer func(storeFile *os.File) {
		err := storeFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(storeFile)

	jsonFile, err := os.Create(migrationsDir + "/migrate.json")
	if err != nil {
		return err
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(jsonFile)

	migrationRunTemplate := template.Must(template.New("main").Parse(string(tpl.InitMigrationRunTemplate())))
	err = migrationRunTemplate.Execute(migrationRunFile, map[string]interface{}{
		"pwd": pwd,
		"dir": migrationsDir,
	})

	if err != nil {
		return err
	}

	storeTemplate := template.Must(template.New("store").Parse(string(tpl.StoreTemplate())))
	err = storeTemplate.Execute(storeFile, nil)
	if err != nil {
		return err
	}

	cmd := exec.Command("gofmt", "-w", migrationsDir+"/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	cmd = exec.Command("gofmt", "-w", migrationsDir+"/store.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}
