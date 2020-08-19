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

package clean

import (
	"fmt"

	"github.com/amarjeetanandsingh/tgconst/files"
)

type cleaner struct {
	// Setting `isRecursive` true will generate const for all the subdirectories too.
	// By default it generates const only for the current directory.
	isRecursive bool

	// Generate const file for the given directory. If `isRecursive` flag is set,
	// it will generate const file recursively for all its subdirectories too.
	dir string

	// print all cleanup operations
	verbose bool
}

func New(options ...func(c *cleaner)) *cleaner {
	c := &cleaner{}

	// set all the options
	for _, option := range options {
		option(c)
	}

	return c
}

// TODO:v2: make _tgconst_gen.go suffix as config
func (c *cleaner) Do() error {
	f := files.New(c.verbose, c.isRecursive, c.dir)
	if err := f.DeleteFilesWithSuffix("_tgconst_gen.go"); err != nil {
		return fmt.Errorf("error deleting %s suffix files in %s dir :: %w", "_tgconst_gen.go", c.dir, err)
	}
	return nil
}
