package migrate

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/fatih/color"
)

func Run(until, migrationsDir string, lastRunTS int, up bool) string {
	migrationsBuild := exec.Command("go", "build")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationsBuild.Dir = wd + "/" + migrationsDir

	errBuild := migrationsBuild.Start()
	if errBuild != nil {
		log.Fatal(errBuild)
	}

	errWait := migrationsBuild.Wait()
	if errWait != nil {
		log.Fatal(errWait)
	}

	var migrationArgs []string

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

	migrationsRun := exec.Command("./"+migrationsDir, migrationArgs...)
	migrationsRun.Dir = wd + "/" + migrationsDir
	var errBuf bytes.Buffer
	migrationsRun.Stderr = &errBuf

	stdout, err := migrationsRun.StdoutPipe()
	_ = migrationsRun.Start()

	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	for err == nil {
		color.Cyan("%v", line)
		line, err = reader.ReadString('\n')
	}

	return errBuf.String()
}
