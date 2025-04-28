package migrate

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/g14a/metana/pkg/store"
	"github.com/spf13/afero"
)

func Run(opts MigrationOptions) (string, error) {
	var scriptsDir string
	if filepath.IsAbs(opts.MigrationsDir) {
		scriptsDir = filepath.Join(opts.MigrationsDir, "scripts")
	} else {
		scriptsDir = filepath.Join(opts.Wd, opts.MigrationsDir, "scripts")
	}

	files, err := filepath.Glob(filepath.Join(scriptsDir, "*.go"))
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no migrations found in %s", scriptsDir)
	}

	sort.Strings(files)

	executed := make(map[string]bool)
	if !opts.DryRun {
		sh, err := store.GetStoreViaConn(opts.StoreConn, opts.MigrationsDir, afero.NewOsFs(), opts.Wd)
		if err == nil {
			track, err := sh.Load(afero.NewOsFs())
			if err == nil {
				for _, m := range track.Migrations {
					executed[m.Title] = true
				}
			}
		}
	}

	var allOutput strings.Builder

	for _, file := range files {
		base := filepath.Base(file)
		migrationName := strings.TrimSuffix(strings.SplitN(base, "_", 2)[1], ".go")

		// In non-dry run mode, skip migrations based on idempotency
		if !opts.DryRun {
			if opts.Up && executed[base] {
				continue
			}
			if !opts.Up && !executed[base] {
				continue
			}
		}

		runMigration := func(mode string) error {
			args := []string{"run", file, "-mode", mode}
			cmd := exec.Command("go", args...)
			cmd.Dir = opts.Wd

			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			stdoutPipe, err := cmd.StdoutPipe()
			if err != nil {
				return fmt.Errorf("stdout error: %w", err)
			}

			if err := cmd.Start(); err != nil {
				return fmt.Errorf("start error: %w", err)
			}

			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				color.Cyan("%s", line)
				allOutput.WriteString(line + "\n")
			}

			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("execution error: %v\n%s", err, stderr.String())
			}
			return nil
		}

		if opts.Up {
			upErr := runMigration("up")
			if upErr != nil {
				color.Red("Migration %s failed, attempting rollback...\n", base)
				downErr := runMigration("down")
				if downErr != nil {
					return allOutput.String(), fmt.Errorf("rollback of %s also failed: %v", base, downErr)
				}
				return allOutput.String(), fmt.Errorf("migration %s failed: %v", base, upErr)
			}
		} else {
			if err := runMigration("down"); err != nil {
				return allOutput.String(), fmt.Errorf("migration %s failed: %v", base, err)
			}
		}

		if opts.Until != "" && migrationName == opts.Until {
			color.Yellow(" >>> Reached --until: %s. Stopping further migrations.\n", opts.Until)
			break
		}
	}

	return allOutput.String(), nil
}

type MigrationOptions struct {
	Until         string
	MigrationsDir string
	Wd            string
	LastRunTS     int
	Up            bool
	StoreConn     string
	DryRun        bool
}
