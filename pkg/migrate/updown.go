package migrate

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
)

func RunUp(until, migrationsDir string) (string, error) {
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

	migrationsRun := exec.Command("./"+migrationsDir, migrationArgs...)
	migrationsRun.Dir = wd + "/" + migrationsDir
	var outBuf, errBuf bytes.Buffer
	migrationsRun.Stdout = &outBuf
	migrationsRun.Stderr = &errBuf

	errRun := migrationsRun.Run()
	if errRun != nil {
		return outBuf.String(), errRun
	}

	if errBuf.Len() > 0 {
		return outBuf.String(), errors.New(errBuf.String())
	}

	return outBuf.String(), nil
}

func RunDown(until, migrationsDir string) (string, error) {
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

	migrationsRun := exec.Command("./"+migrationsDir, migrationArgs...)
	migrationsRun.Dir = wd + "/" + migrationsDir
	var outBuf, errBuf bytes.Buffer
	migrationsRun.Stdout = &outBuf
	migrationsRun.Stderr = &errBuf

	errRun := migrationsRun.Run()
	if errRun != nil {
		return outBuf.String(), errRun
	}

	if errBuf.Len() > 0 {
		return outBuf.String(), errors.New(errBuf.String())
	}

	return outBuf.String(), nil
}
