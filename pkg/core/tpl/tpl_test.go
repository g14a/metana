package tpl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandaloneMigrationTemplate_Default(t *testing.T) {
	result := StandaloneMigrationTemplate("", "")
	content := string(result)

	assert.True(t, strings.Contains(content, "func up() error {"))
	assert.True(t, strings.Contains(content, `fmt.Println("{{ .MigrationName }} up")`))
	assert.True(t, strings.Contains(content, "func down() error {"))
	assert.True(t, strings.Contains(content, `fmt.Println("{{ .MigrationName }} down")`))
	assert.True(t, strings.Contains(content, "flag.String(\"mode\", \"up\""))
	assert.True(t, strings.Contains(content, "{{ .Filename }}"))
}

func TestStandaloneMigrationTemplate_Custom(t *testing.T) {
	customUp := `fmt.Println("custom up")`
	customDown := `fmt.Println("custom down")`

	result := StandaloneMigrationTemplate(customUp, customDown)
	content := string(result)

	assert.True(t, strings.Contains(content, customUp))
	assert.True(t, strings.Contains(content, customDown))
	assert.True(t, strings.Contains(content, "func up() error {"))
	assert.True(t, strings.Contains(content, "func down() error {"))
}
