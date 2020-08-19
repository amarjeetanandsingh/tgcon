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

	"github.com/amarjeetanandsingh/tgconst/config"
	"github.com/amarjeetanandsingh/tgconst/gen"

	"github.com/spf13/cobra"
)

const (
	genShortDoc = "generates struct field tag values as string constant"
	genLongDoc  = `
It generates string constants of struct field tag values. All constants are
generated into a new file(with _tgconst_gen.go as suffix) for each package.`
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: genShortDoc,
	Long:  genLongDoc,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetGeneratorCfg()
		g := gen.New(
			gen.Dir(cfg.Dir),
			gen.Tags(cfg.Tags),
			gen.Recursive(cfg.IsRecursive),
			gen.TaggedFieldOnly(cfg.OnlyTaggedFields),
			gen.MissingTagPolicy(cfg.MissingTagValFormat),
		)

		if err := g.Do(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	cfg := config.GetGeneratorCfg()
	genCmd.Flags().StringVarP(&cfg.Dir, "dir", "d", ".", "Generate tag as constants for the given dir directory")
	genCmd.Flags().StringSliceVarP(&cfg.Tags, "tags", "t", []string{}, "Comma separated list of tags we are going to create constants for")
	genCmd.Flags().BoolVarP(&cfg.IsRecursive, "recursive", "r", false, "Recursively create constants for all subdirectories too")
	genCmd.Flags().BoolVarP(&cfg.OnlyTaggedFields, "onlyTagged", "e", false, "Escape empty tag fields. Do not create string constants for unTagged fields.")
	genCmd.Flags().StringVarP(&cfg.MissingTagValFormat, "missingTagValFormat", "f", "", "policy to generate tag value for fields with no tags. Currently supports [camelcase, lispcase, pascalcase, snakecase]")
}
