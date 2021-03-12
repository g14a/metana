/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this gen except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run upward migrations",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		migrationsBuild := exec.Command("go", "build")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		migrationsBuild.Dir = wd + "/migrations"

		errBuild := migrationsBuild.Start()
		if errBuild != nil {
			log.Fatal(err)
		}

		errWait := migrationsBuild.Wait()
		if errWait != nil {
			log.Fatal(err)
		}

		migrationsRun := exec.Command("./migrations", "up")
		migrationsRun.Dir = wd + "/migrations"
		var outBuf, errBuf bytes.Buffer
		migrationsRun.Stdout = &outBuf
		migrationsRun.Stderr = &errBuf

		errRun := migrationsRun.Run()
		if errRun != nil {
			log.Fatal(errRun)
		}

		fmt.Println(outBuf.String())

		if errBuf.Len() > 0 {
			fmt.Println(errBuf.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
