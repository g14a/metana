package migrate

import (
	"testing"
	"time"

	"github.com/g14a/metana/pkg"
)

func TestRun(t *testing.T) {
	// Run full up
	logs := Run("", "testdata", "../../..", 0, true)
	pkg.ExpectLines(t, logs, []string{"1621746399-InitSchema.go", "1621746406-AddIndexes.go", "1621746410-AddData.go"}...)

	// Run full down
	logs = Run("", "testdata", "../../..", int(time.Now().Unix()), false)
	pkg.ExpectLines(t, logs, []string{"1621746410-AddData.go", "1621746406-AddIndexes.go", "1621746399-InitSchema.go"}...)

	// Run up Until
	logs = Run("AddIndexes", "testdata", "../../..", 0, true)
	pkg.ExpectLines(t, logs, []string{"1621746399-InitSchema.go", "1621746406-AddIndexes.go"}...)

	// Run Down Until
	logs = Run("AddIndexes", "testdata", "../../..", int(time.Now().Unix()), false)
	pkg.ExpectLines(t, logs, []string{"1621746406-AddIndexes.go", "1621746406-AddIndexes.go"}...)
}
