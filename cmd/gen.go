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

	"github.com/amarjeetanandsingh/tgcon/config"
	"github.com/amarjeetanandsingh/tgcon/gen"

	"github.com/spf13/cobra"
)

const (
	genShortDoc = "generates struct field tag values as string constant"
	genLongDoc  = `
It generates string constants of struct field tag values. All constants are
generated into a new file(with _tgcon_gen.go as suffix) for each package.`
)

func NewGenCmd() *cobra.Command {
	genCmd := &cobra.Command{
		Use:   "gen",
		Short: genShortDoc,
		Long:  genLongDoc,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetGeneratorCfg()
			generator := gen.New(
				gen.All(cfg.All),
				gen.Dir(cfg.Dir),
				gen.Tags(cfg.Tags),
				gen.Recursive(cfg.IsRecursive),
				gen.TaggedFieldOnly(cfg.OnlyTaggedFields),
				gen.DefaultTagValFormat(cfg.DefaultTagValFormat),
			)

			if err := generator.Do(); err != nil {
				return err
			}
			fmt.Println("Done.")
			return nil
		},
	}

	// set flags
	setGenFlags(genCmd)

	return genCmd
}

func setGenFlags(genCmd *cobra.Command) {
	cfg := config.GetGeneratorCfg()
	genCmd.Flags().BoolVarP(&cfg.All, "all", "a", false, "Process all structs irrespective of magic comment")
	genCmd.Flags().StringVarP(&cfg.Dir, "dir", "d", ".", "Generate constants for the given dir")
	genCmd.Flags().StringSliceVarP(&cfg.Tags, "tags", "t", []string{}, "Create constants only for given comma separated list of tags. Empty means process all available tags")
	genCmd.Flags().BoolVarP(&cfg.IsRecursive, "recursive", "r", false, "Recursively create constants for all subdirectories too")
	genCmd.Flags().BoolVarP(&cfg.OnlyTaggedFields, "onlyTagged", "s", false, "Do not create constants for unTagged fields. -s means skip")
	genCmd.Flags().StringVarP(&cfg.DefaultTagValFormat, "defaultTagValFormat", "f", "", "Format to generate tag value constant for fields with no tags. Default format is  Currently supports [camelcase, lispcase, pascalcase, snakecase]")
}
