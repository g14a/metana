/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
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
	"io/ioutil"
	"log"
	"os"

	"github.com/g14a/metana/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage your local metana config in .metana.yml",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var setConfigCmd = &cobra.Command{
	Use:   "set",
	Short: "Set config",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatal(err)
		}

		store, err := cmd.Flags().GetString("store")
		if err != nil {
			log.Fatal(err)
		}

		mc, err := config.GetMetanaConfig()
		if os.IsNotExist(err) {
			_, err = os.Create(".metana.yml")
			if err != nil {
				log.Fatal(err)
			}
		}

		if dir != "" {
			mc.Dir = dir
		}

		if store != "" {
			mc.StoreConn = store
		}

		b, err := yaml.Marshal(mc)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(".metana.yml", b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	setConfigCmd.Flags().StringP("store", "s", "", "Set your store")
	setConfigCmd.Flags().StringP("dir", "d", "migrations", "Set your migrations directory")
	configCmd.AddCommand(setConfigCmd)
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
