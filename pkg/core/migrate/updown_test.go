package migrate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/g14a/metana/pkg"
)

func TestRun(t *testing.T) {
	// Run full up
	opts := MigrationOptions{
		Until:         "",
		MigrationsDir: "testdata",
		Wd:            "../../..",
		LastRunTS:     0,
		Up:            true,
	}
	logs, err := Run(opts)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746399-InitSchema.go", "1621746406-AddIndexes.go", "1621746410-AddData.go"}...)

	opts.LastRunTS = int(time.Now().Unix())
	opts.Up = false

	// Run full down
	logs, err = Run(opts)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746410-AddData.go", "1621746406-AddIndexes.go", "1621746399-InitSchema.go"}...)

	// Run up Until
	opts.Until = "AddIndexes"
	opts.LastRunTS = 0
	opts.Up = true

	logs, err = Run(opts)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746399-InitSchema.go", "1621746406-AddIndexes.go"}...)

	// Run Down Until
	opts.Until = "AddIndexes"
	opts.LastRunTS = int(time.Now().Unix())
	opts.Up = false

	logs, err = Run(opts)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746406-AddIndexes.go", "1621746406-AddIndexes.go"}...)
}
