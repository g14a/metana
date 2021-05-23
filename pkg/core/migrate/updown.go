package migrate

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strconv"

	"github.com/fatih/color"
)

func Run(until, migrationsDir string, wd string, lastRunTS int, up bool) string {
	var migrationArgs []string

	migrationArgs = append(migrationArgs, "run", "main.go")

	if up {
		migrationArgs = append(migrationArgs, "up")
	} else {
		migrationArgs = append(migrationArgs, "down")
	}

	if until != "" {
		migrationArgs = append(migrationArgs, "--until", until)
	}
	lastRunTSString := strconv.Itoa(lastRunTS)

	migrationArgs = append(migrationArgs, "--last-run-ts", lastRunTSString)

	migrationsRun := exec.Command("go", migrationArgs...)
	migrationsRun.Dir = wd + "/" + migrationsDir
	var errBuf bytes.Buffer
	migrationsRun.Stderr = &errBuf

	stdout, err := migrationsRun.StdoutPipe()
	err = migrationsRun.Start()
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		color.Cyan("%v", line)
		line, err = reader.ReadString('\n')
	}

	return errBuf.String()
}
