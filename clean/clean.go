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
	"io"
	"os"
	"path"
	"strings"

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

	// write to stdOut
	stdOut io.Writer
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
	if err := c.deleteFilesWithSuffix(c.dir, "_tgconst_gen.go"); err != nil {
		return fmt.Errorf("error deleting %s suffix files in %s dir :: %w", "_tgconst_gen.go", c.dir, err)
	}
	return nil
}

func (c *cleaner) deleteFilesWithSuffix(dir, suffix string) error {
	fileNames, err := files.ListFilesInDir(dir, func(fileName string) bool {
		return strings.HasSuffix(fileName, suffix)
	})
	if err != nil {
		return fmt.Errorf("files.ListFilesInDir():: error getting files. dir = %s, suffix=%s :: %w", dir, suffix, err)
	}

	// delete files in current directory
	for _, fileName := range fileNames {
		filePath := path.Join(dir, fileName)
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("s.Remove(filePath):: error deleting %s file :: %w", filePath, err)
		}

		// log
		if c.verbose {
			fmt.Fprintln(c.stdOut, "Deleted:", filePath)
		}
	}

	// Stop processing if it's not recursive
	if !c.isRecursive {
		return nil
	}

	// ******* process recursive call *********

	subDirs, err := files.ListDirs(dir, func(dirName string) bool {
		// ignore hidden dirs
		return len(dirName) > 0 && dirName[0] != '.'
	})
	if err != nil {
		return fmt.Errorf("files.ListDirs():: error getting dirs. dir = %s :: %w", dir, err)
	}

	for _, subDir := range subDirs {
		subDirPath := path.Join(dir, subDir)
		if err := c.deleteFilesWithSuffix(subDirPath, suffix); err != nil {
			return err
		}
	}
	return nil
}
