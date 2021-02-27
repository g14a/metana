package tpl

func MainTemplate() []byte {
	return []byte(`
import (
	"fmt"
)

type {{ .Name }} struct {

}

func (r *{{ .Name }}) Up()  {
	fmt.Println("probation up")
}

func (r *{{ .Name }}) Down()  {
	fmt.Println("probation down")
}
`)
}
