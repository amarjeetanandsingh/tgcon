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

package files

import (
	"fmt"
	"io/ioutil"
)

// ListFilesInDir returns the list of all files inside the given dir, satisfying
// checkName func. Nil checkName func allows all files. It's not recursive.
func ListFilesInDir(dir string, checkName func(string) bool) ([]string, error) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("files.ListFilesInDir() :: error reading path: %s :: %w", dir, err)
	}

	var fileNames []string

	for _, f := range ff {
		if !f.IsDir() && (checkName == nil || checkName(f.Name())) {
			fileNames = append(fileNames, f.Name())
		}
	}
	return fileNames, nil
}

// ListDirs returns all the immediate sub-directory names inside given dir satisfying
// checkName func. If checkName func is nil, it returns all immediate sub-directories.
// This is not recursive.
func ListDirs(dir string, checkName func(string) bool) ([]string, error) {
	ff, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("files.listDirs():: error reading path: %s :: %w", dir, err)
	}

	var dirNames []string

	for _, f := range ff {
		if f.IsDir() && (checkName == nil || checkName(f.Name())) {
			dirNames = append(dirNames, f.Name())
		}
	}
	return dirNames, nil
}
