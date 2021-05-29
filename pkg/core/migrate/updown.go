package migrate

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/spf13/cobra"

	"github.com/fatih/color"
)

func Run(opts MigrationOptions) (string, error) {
	migrationArgs := GetMigrationArgs(opts)

	var envKeys []string
	if opts.EnvFile != "" {
		envMap, err := godotenv.Read(opts.Wd + "/" + opts.EnvFile)
		if err != nil {
			return "", err
		}
		for k, v := range envMap {
			envKeys = append(envKeys, fmt.Sprintf("%s=%s", k, v))
		}
	}

	fmt.Println(migrationArgs, "=========migration args========")
	fmt.Println(opts.MigrationsDir, "=========mig dir==========")

	migrationsRun := exec.Command("go", migrationArgs...)
	migrationsRun.Env = append(os.Environ(), envKeys...)
	migrationsRun.Dir = opts.Wd + "/" + opts.MigrationsDir

	var errBuf bytes.Buffer
	migrationsRun.Stderr = &errBuf

	stdout, err := migrationsRun.StdoutPipe()
	if err != nil {
		return "", err
	}
	err = migrationsRun.Start()
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		color.Cyan("%v", line)
		line, err = reader.ReadString('\n')
	}

	return errBuf.String(), nil
}

func GetMigrationArgs(opts MigrationOptions) []string {
	var migrationArgs []string

	migrationArgs = append(migrationArgs, "run", "main.go")

	if opts.Up {
		migrationArgs = append(migrationArgs, "up")
	} else {
		migrationArgs = append(migrationArgs, "down")
	}

	if opts.Until != "" {
		migrationArgs = append(migrationArgs, "--until", opts.Until)
	}

	lastRunTSString := strconv.Itoa(opts.LastRunTS)
	migrationArgs = append(migrationArgs, "--last-run-ts", lastRunTSString)

	return migrationArgs
}

type MigrationOptions struct {
	Until         string
	MigrationsDir string
	Wd            string
	LastRunTS     int
	Up            bool
	StoreConn     string
	DryRun        bool
	EnvFile       string
	Cmd           *cobra.Command
}
