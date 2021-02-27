package gen

import (
	"fmt"
	"go-migrate/pkg/tpl"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
)

func CreateMigrationFile(file string) (string, error) {
	nm := tpl.NewMigration{
		Name:      strcase.ToCamel(file),
		Timestamp: strconv.Itoa(int(time.Now().Unix())),
	}

	fileName := fmt.Sprintf("migrations/%s-%s.go", nm.Timestamp, nm.Name)

	mainFile, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MigrationTemplate())))
	err = mainTemplate.Execute(mainFile, nm)
	if err != nil {
		return "", err
	}

	migrationRunFile, err := os.Create("migrations/main.go")
	migrationRunTemplate := template.Must(template.New("main").Parse(string(tpl.InitMigrationRunTemplate())))
	err = migrationRunTemplate.Execute(migrationRunFile, nil)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func CreateInitConfig() {
	initInterface, err := os.Create("migrations/interfaces/interface.go")
	if err != nil {
		fmt.Println(err)
	}

	defer initInterface.Close()

	mainTemplate := template.Must(template.New("init").Parse(string(tpl.InitMigrationTemplate())))
	err = mainTemplate.Execute(initInterface, nil)
	if err != nil {
		fmt.Println(err)
	}
}
