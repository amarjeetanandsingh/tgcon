/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"

	"github.com/amarjeetanandsingh/tgconst/clean"
	"github.com/amarjeetanandsingh/tgconst/config"

	"github.com/spf13/cobra"
)

const (
	// TODO
	cleanShortDoc = ""
	cleanLongDoc  = ""
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Long:  cleanLongDoc,
	Short: cleanShortDoc,
	Run: func(cmd *cobra.Command, args []string) {
		c := clean.New(
			clean.Dir(config.GetCleanerCfg().Dir),
			clean.Verbose(config.GetCleanerCfg().Verbose),
			clean.Recursive(config.GetCleanerCfg().IsRecursive),
		)
		if err := c.Do(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cfg := config.GetCleanerCfg()
	cleanCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "verbose output")
	cleanCmd.Flags().StringVarP(&cfg.Dir, "dir", "d", ".", "Delete generated const file from given directory")
	cleanCmd.Flags().BoolVarP(&cfg.IsRecursive, "recursive", "r", false, "Recursively delete generated const files for all subdirectories too")
}
