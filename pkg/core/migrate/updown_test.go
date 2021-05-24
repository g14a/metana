package migrate

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/g14a/metana/pkg"
)

func TestRun(t *testing.T) {
	// Run full up
	logs, err := Run("", "testdata", "../../..", 0, true)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746399-InitSchema.go", "1621746406-AddIndexes.go", "1621746410-AddData.go"}...)

	// Run full down
	logs, err = Run("", "testdata", "../../..", int(time.Now().Unix()), false)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746410-AddData.go", "1621746406-AddIndexes.go", "1621746399-InitSchema.go"}...)

	// Run up Until
	logs, err = Run("AddIndexes", "testdata", "../../..", 0, true)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746399-InitSchema.go", "1621746406-AddIndexes.go"}...)

	// Run Down Until
	logs, err = Run("AddIndexes", "testdata", "../../..", int(time.Now().Unix()), false)
	assert.NoError(t, err)
	pkg.ExpectLines(t, logs, []string{"1621746406-AddIndexes.go", "1621746406-AddIndexes.go"}...)
}
