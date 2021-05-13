package migrate

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func RunUp(until, migrationsDir string, lastRunTS int) (string, string) {
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

	migrationArgs := []string{"up"}
	if until != "" {
		migrationArgs = append(migrationArgs, "--until", until)
	}
	lastRunTSString := strconv.Itoa(lastRunTS)

	migrationArgs = append(migrationArgs, "--last-run-ts", lastRunTSString)

	migrationsRun := exec.Command("./"+migrationsDir, migrationArgs...)
	migrationsRun.Dir = wd + "/" + migrationsDir
	var outBuf, errBuf bytes.Buffer
	// migrationsRun.Stdout = &outBuf
	migrationsRun.Stderr = &errBuf

	stdout, err := migrationsRun.StdoutPipe()
	migrationsRun.Start()

	oneByte := make([]byte, 256)

	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			break
		}
		fmt.Println(string(oneByte), "============one byte=========")
	}

	return outBuf.String(), errBuf.String()
}

func RunDown(until, migrationsDir string, lastRunTS int) (string, string) {
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

	migrationArgs := []string{"down"}
	if until != "" {
		migrationArgs = append(migrationArgs, "--until", until)
	}

	lastRunTSString := strconv.Itoa(lastRunTS)
	migrationArgs = append(migrationArgs, "--last-run-ts", lastRunTSString)

	migrationsRun := exec.Command("./"+migrationsDir, migrationArgs...)
	migrationsRun.Dir = wd + "/" + migrationsDir
	var outBuf, errBuf bytes.Buffer
	migrationsRun.Stdout = &outBuf
	migrationsRun.Stderr = &errBuf

	errRun := migrationsRun.Run()
	if errRun != nil {
		return outBuf.String(), errRun.Error()
	}

	return outBuf.String(), errBuf.String()
}
