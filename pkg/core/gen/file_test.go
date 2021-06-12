package gen

import (
	"os"
	"testing"

	"github.com/g14a/metana/pkg"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestCreateInitConfig(t *testing.T) {
	FS := afero.NewMemMapFs()

	os.Chdir("/Users/g14a/metana")

	FS.MkdirAll("/Users/g14a/metana/migrations", 0755)
	FS.MkdirAll("/Users/g14a/metana/migrations/scripts", 0755)

	err := CreateInitConfig("migrations", "github.com/g14a/metana", FS, "")
	assert.Equal(t, true, err == nil)

	bytes, err := afero.ReadFile(FS, "migrations/main.go")
	assert.Equal(t, true, err == nil)

	expectedLines := getExpectedLinesInit()

	pkg.ExpectLines(t, string(bytes), expectedLines...)
}

func TestMigrationExists(t *testing.T) {

	FS := afero.NewMemMapFs()

	FS.MkdirAll("/Users/g14a/metana/migrations/scripts", 0755)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621081055-InitSchema.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621084125-AddIndexes.go", []byte("{}"), 0644)
	afero.WriteFile(FS, "/Users/g14a/metana/migrations/scripts/1621084135-AddFKeys.go", []byte("{}"), 0644)

	os.Chdir("/Users/g14a/metana")

	tests := []struct {
		inputMigrationsDir string
		inputMigrationName string
		FS                 afero.Fs
		Exists             bool
		Environment        string
	}{
		{
			inputMigrationsDir: "migrations",
			inputMigrationName: "addFKeys",
			FS:                 FS,
			Exists:             true,
			Environment:        "",
		}, {
			inputMigrationsDir: "migrations",
			inputMigrationName: "FKeys",
			FS:                 FS,
			Exists:             false,
			Environment:        "",
		}, {
			inputMigrationsDir: "migrations",
			inputMigrationName: "AddFKeys",
			FS:                 FS,
			Exists:             true,
			Environment:        "",
		}, {
			inputMigrationsDir: "migrations",
			inputMigrationName: "initSchema",
			FS:                 FS,
			Exists:             true,
			Environment:        "",
		}, {
			inputMigrationsDir: "migrations",
			inputMigrationName: "addIndexes",
			FS:                 FS,
			Exists:             true,
			Environment:        "",
		}, {
			inputMigrationsDir: "migrations",
			inputMigrationName: "nsdlgvnw",
			FS:                 FS,
			Exists:             false,
			Environment:        "",
		},
	}

	for _, tt := range tests {
		exists, err := MigrationExists("/Users/g14a/metana", tt.inputMigrationsDir, tt.inputMigrationName, tt.FS, tt.Environment)
		assert.Equal(t, tt.Exists, exists)
		assert.Equal(t, true, err == nil)
	}
}

func TestCreateMigrationFile(t *testing.T) {
	FS := afero.NewMemMapFs()

	os.Chdir("/Users/g14a/metana")
	FS.MkdirAll("/Users/g14a/metana/migrations", 0755)

	cOpts := CreateMigrationOpts{
		Wd:            "/Users/g14a/metana",
		MigrationsDir: "migrations",
		File:          "initSchema",
		CustomTmpl:    "",
		Environment:   "",
		FS:            FS,
	}

	filename, err := CreateMigrationFile(cOpts)
	assert.Equal(t, true, err == nil)

	resultFileBytes, err := afero.ReadFile(FS, filename)
	assert.Equal(t, true, err == nil)

	pkg.ExpectLines(t, string(resultFileBytes), getExpectedLinesMigration()...)
}

func getExpectedLinesMigration() []string {
	return []string{`type InitSchemaMigration struct`,
		`Timestamp     int`,
		`Filename      string`,
		`MigrationName string`,
		`InitSchemaMigration`,
		`"InitSchema up"`,
		`"InitSchema down"`}
}

func getExpectedLinesInit() []string {
	return []string{`// This file is auto generated. DO NOT EDIT!`,
		`MigrateUp`,
		`upUntil string, lastRunTS int`,
		`return nil`,
		`MigrateDown`,
		`downUntil string, lastRunTS int`,
		`return nil`,
		`func main()`,
		`upCmd := flag.NewFlagSet`,
		`downCmd := flag.NewFlagSet`,
		`var upUntil, downUntil string`,
		`var lastRunTS int`,
		`upCmd.StringVar`,
		`upCmd.IntVar`,
		`downCmd.StringVar`,
		`downCmd.IntVar`,
		`switch`,
		`case "up"`,
		`err := upCmd.Parse`,
		`if err != nil {`,
		`return`,
		`}`,
		`case "down"`,
		`err := downCmd.Parse`,
		`if err != nil {`,
		`return`,
		`}`,
		`if upCmd.Parsed()`,
		`MigrateUp()`,
		`if err != nil {`,
		`fmt.Fprintf()`,
		`}`,
		`if downCmd.Parsed()`,
		`MigrateDown()`,
		`if err != nil {`,
		`fmt.Fprintf()`,
		`}`,
	}
}
