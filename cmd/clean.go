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

	"github.com/amarjeetanandsingh/tgcon/clean"
	"github.com/amarjeetanandsingh/tgcon/config"

	"github.com/spf13/cobra"
)

const (
	// TODO
	cleanShortDoc = ""
	cleanLongDoc  = ""
)

// cleanCmd represents the clean sub-command
func NewCleanCmd() *cobra.Command {
	cleanCmd := &cobra.Command{
		Use:   "clean",
		Long:  cleanLongDoc,
		Short: cleanShortDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			cleaner := clean.New(
				clean.StdOut(cmd.OutOrStdout()),
				clean.Dir(config.GetCleanerCfg().Dir),
				clean.Verbose(config.GetCleanerCfg().Verbose),
				clean.Recursive(config.GetCleanerCfg().IsRecursive),
			)
			if err := cleaner.Do(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Done.")
			return nil
		},
	}

	setCleanFlags(cleanCmd)

	return cleanCmd
}

func setCleanFlags(cleanCmd *cobra.Command) {
	cfg := config.GetCleanerCfg()
	cleanCmd.Flags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "verbose output")
	cleanCmd.Flags().StringVarP(&cfg.Dir, "dir", "d", ".", "Delete generated const file from given directory")
	cleanCmd.Flags().BoolVarP(&cfg.IsRecursive, "recursive", "r", false, "Recursively delete generated const files for all subdirectories too")
}
