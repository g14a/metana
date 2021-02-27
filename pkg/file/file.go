package file

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"go-migrate/pkg/tpl"
	"os"
	"strconv"
	"text/template"
	"time"
)

func CreateMigrationFile(file string) {
	nm := tpl.NewMigration{
		Name:      strcase.ToCamel(file),
		Timestamp: strconv.Itoa(int(time.Now().Unix())),
	}

	mainFile, err := os.Create(fmt.Sprintf("migrations/%s-%s.go", nm.Timestamp, nm.Name))
	if err != nil {
		fmt.Println(err)
	}

	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, nm)
	if err != nil {
		fmt.Println(err)
	}
}

func CreateInitConfig() {
	initInterface, err := os.Create("migrations/interfaces/interface.go")
	if err != nil {
		fmt.Println(err)
	}

	defer initInterface.Close()

	mainTemplate := template.Must(template.New("init").Parse(string(tpl.InitTemplate())))
	err = mainTemplate.Execute(initInterface, nil)
	if err != nil {
		fmt.Println(err)
	}
}
