package tpl

func MainTemplate() []byte {
	return []byte(`
import (
	"fmt"
)

type {{ .Name }} struct {

}

func (r *{{ .Name }}) Up()  {
	fmt.Println("migration up")
}

func (r *{{ .Name }}) Down()  {
	fmt.Println("migration down")
}
`)
}

func InitTemplate() []byte {
	return []byte(`package _interface

type Migration interface {
	Up()
	Down()
}
`)
}

type NewMigration struct {
	Name      string
	Timestamp string
}
