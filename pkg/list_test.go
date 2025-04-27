package pkg

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetMigrations_validFiles(t *testing.T) {
	FS := afero.NewMemMapFs()
	base := t.TempDir()
	migrationsDir := filepath.Join(base, "migrations", "scripts")
	_ = FS.MkdirAll(migrationsDir, 0755)

	files := []string{
		"1621081055-InitSchema.go",
		"1621084125-AddIndexes.go",
		"1621084135-AddFKeys.go",
	}
	for _, name := range files {
		afero.WriteFile(FS, filepath.Join(migrationsDir, name), []byte("{}"), 0644)
	}

	migrations, err := GetMigrations(base, "migrations", FS)
	assert.NoError(t, err)

	for i, m := range migrations {
		assert.Equal(t, files[i], m.Name)
	}
}

func TestGetMigrations_no_files(t *testing.T) {
	FS := afero.NewMemMapFs()
	base := t.TempDir()
	_ = FS.MkdirAll(filepath.Join(base, "migrations", "scripts"), 0755)

	migrations, err := GetMigrations(base, "migrations", FS)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(migrations))
}

func TestExecutedAtTracking(t *testing.T) {
	FS := afero.NewMemMapFs()
	base := t.TempDir()
	scripts := filepath.Join(base, "migrations", "scripts")
	_ = FS.MkdirAll(scripts, 0755)

	files := []string{
		"1621081055-InitSchema.go",
		"1621084125-AddIndexes.go",
		"1621084135-AddFKeys.go",
	}
	for _, name := range files {
		afero.WriteFile(FS, filepath.Join(scripts, name), []byte("// mock"), 0644)
	}

	// Before migration execution: executed_at should not exist
	migrations, err := GetMigrations(base, "migrations", FS)
	assert.NoError(t, err)
	for _, m := range migrations {
		assert.Empty(t, m.ModTime)
	}

	// Simulate store after execution
	now := time.Now().Format("02-01-2006 15:04")
	mockStore := &MockStore{
		Data: map[string]string{
			files[0]: now,
			files[1]: now,
			files[2]: now,
		},
	}

	executed := map[string]string{}
	track, err := mockStore.Load(FS)
	assert.NoError(t, err)
	for _, m := range track.Migrations {
		executed[m.Title] = m.ExecutedAt
	}

	for _, name := range files {
		assert.Equal(t, now, executed[name])
	}
}
