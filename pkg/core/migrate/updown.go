package migrate

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/fatih/color"
)

func Run(opts MigrationOptions) (string, error) {
	migrationArgs := GetMigrationArgs(opts)

	var envKeys []string
	if opts.EnvFile != "" {
		_, err := os.Stat(opts.Wd + "/" + opts.EnvFile)
		if !os.IsNotExist(err) {
			envMap, err := godotenv.Read(opts.Wd + "/" + opts.EnvFile)
			if err != nil {
				return "", err
			}
			for k, v := range envMap {
				envKeys = append(envKeys, fmt.Sprintf("%s=%s", k, v))
			}
		}
	}

	fmt.Println(opts.Wd + "/" + opts.MigrationsDir + "/environments/" + opts.Environment)
	migrationsRun := exec.Command("go", migrationArgs...)
	migrationsRun.Env = append(os.Environ(), envKeys...)
	if opts.Environment == "" {
		migrationsRun.Dir = opts.Wd + "/" + opts.MigrationsDir
	} else {
		migrationsRun.Dir = opts.Wd + "/" + opts.MigrationsDir + "/environments/" + opts.Environment
	}

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
	Environment   string
}
