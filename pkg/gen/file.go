package gen

import (
	"fmt"
	"go-migrate/pkg/tpl"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func CreateInitConfig() error {

	migrationRunFile, err := os.Create("migrations/main.go")
	if err != nil {
		return err
	}

	defer migrationRunFile.Close()

	migrationRunTemplate := template.Must(template.New("main").Parse(string(tpl.InitMigrationRunTemplate())))
	err = migrationRunTemplate.Execute(migrationRunFile, nil)
	if err != nil {
		return err
	}

	cmd := exec.Command("gofmt", "-w", "migrations/main.go")
	if errOut, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, ""), err, errOut))
	}

	return nil
}
