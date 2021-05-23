package core

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestParseCustomTemplate(t *testing.T) {
	FS := afero.NewMemMapFs()

	afero.WriteFile(FS, "/Users/g14a/metana/tmpl.go", getCustomTmplContent(), 0644)

	up, down := ParseCustomTemplate("/Users/g14a/metana", "tmpl.go", FS)

	assert.Equal(t, "fmt.Println(\"template up\")\n\treturn nil", up)
	assert.Equal(t, "fmt.Println(\"template down\")\n\treturn nil", down)
}

func getCustomTmplContent() []byte {
	return []byte(`package main

import "fmt"

func Up() error {
	fmt.Println("template up")
	return nil
}

func Down() error {
	fmt.Println("template down")
	return nil
}
`)
}
