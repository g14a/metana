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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/g14a/go-migrate/pkg/gen"
	"log"
	"os"
	"strings"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a migration script in Go",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		exists, err := gen.MigrationExists(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if exists {
			color.Yellow("Migration already exists")
			os.Exit(0)
		}

		fileName, err := gen.CreateMigrationFile(args[0])
		if err != nil {
			color.Yellow("\nTry initializing migration using `go-migrate init`\n\n")
			os.Exit(0)
		}

		wd, _ := os.Getwd()
		color.Green(" ✓ Created " + wd + "/" + fileName)
		color.Green(" ✓ Generated " + wd + "/migrations/main.go")

		err = gen.AddMigration(args[0], strings.TrimPrefix(fileName, "migrations/"))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
