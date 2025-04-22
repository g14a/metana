package migrate

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"text/template"
	"time"

	"github.com/g14a/metana/pkg/core/tpl"
	"github.com/stretchr/testify/assert"
)

// createMigrationFile generates a real .go migration using the actual standalone template
func createMigrationFile(t *testing.T, dir, name string) string {
	t.Helper()
	filename := fmt.Sprintf("%d_%s.go", time.Now().UnixNano(), name)
	filepath := filepath.Join(dir, "scripts", filename)

	tplData := map[string]string{
		"MigrationName": name,
		"Filename":      filename,
	}

	migTpl := tpl.StandaloneMigrationTemplate("", "")
	tmpl := template.Must(template.New("mig").Parse(string(migTpl)))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, tplData)
	assert.NoError(t, err)

	err = os.WriteFile(filepath, buf.Bytes(), 0644)
	assert.NoError(t, err)

	// Sleep to ensure unique filenames
	time.Sleep(5 * time.Millisecond)
	return filepath
}

// setupTestMigrations creates a temp dir with valid migration scripts
func setupTestMigrations(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	scriptsDir := filepath.Join(dir, "scripts")

	err := os.MkdirAll(scriptsDir, 0755)
	assert.NoError(t, err)

	for _, name := range []string{"InitSchema", "AddIndexes", "AddData"} {
		createMigrationFile(t, dir, name)
	}

	return dir
}

func TestRunUpDown_WithFileStore(t *testing.T) {
	dir := setupTestMigrations(t)
	storeFile := filepath.Join(dir, "migrate.json")
	_ = os.Remove(storeFile)

	opts := MigrationOptions{
		MigrationsDir: ".",
		Wd:            dir,
		StoreConn:     "",
		EnvFile:       "",
	}

	// Run full up
	opts.Up = true
	logs, err := Run(opts)
	assert.NoError(t, err)
	assert.Contains(t, logs, "InitSchema up")
	assert.Contains(t, logs, "AddIndexes up")
	assert.Contains(t, logs, "AddData up")

	// Re-run up: nothing should run
	logs, err = Run(opts)
	assert.NoError(t, err)
	assert.NotContains(t, logs, "InitSchema up")
	assert.NotContains(t, logs, "AddIndexes up")
	assert.NotContains(t, logs, "AddData up")

	// Run full down
	opts.Up = false
	logs, err = Run(opts)
	assert.NoError(t, err)
	assert.Contains(t, logs, "AddData down")
	assert.Contains(t, logs, "AddIndexes down")
	assert.Contains(t, logs, "InitSchema down")

	// Re-run down: nothing should run
	logs, err = Run(opts)
	assert.NoError(t, err)
	assert.NotContains(t, logs, "AddData down")
	assert.NotContains(t, logs, "AddIndexes down")
	assert.NotContains(t, logs, "InitSchema down")
}

func TestRun_UpUntilStop(t *testing.T) {
	dir := setupTestMigrations(t)

	opts := MigrationOptions{
		MigrationsDir: ".",
		Wd:            dir,
		StoreConn:     "",
		Until:         "AddIndexes",
		Up:            true,
	}

	logs, err := Run(opts)
	assert.NoError(t, err)
	assert.Contains(t, logs, "InitSchema up")
	assert.Contains(t, logs, "AddIndexes up")
	assert.NotContains(t, logs, "AddData up")
}

func TestRun_DownUntilStop(t *testing.T) {
	dir := setupTestMigrations(t)

	// Up all
	_, err := Run(MigrationOptions{
		MigrationsDir: ".",
		Wd:            dir,
		StoreConn:     "",
		Up:            true,
	})
	assert.NoError(t, err)

	// Down until AddIndexes
	opts := MigrationOptions{
		MigrationsDir: ".",
		Wd:            dir,
		StoreConn:     "",
		Until:         "AddIndexes",
		Up:            false,
	}

	logs, err := Run(opts)
	assert.NoError(t, err)
	assert.Contains(t, logs, "AddData down")
	assert.Contains(t, logs, "AddIndexes down")
	assert.NotContains(t, logs, "InitSchema down")
}
