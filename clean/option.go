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

import "io"

// moved the option functions from clean.go to keep the logic clean out there

func Dir(d string) func(*cleaner) {
	return func(c *cleaner) {
		c.dir = d
	}
}

func Recursive(isRecursive bool) func(*cleaner) {
	return func(c *cleaner) {
		c.isRecursive = isRecursive
	}
}

func Verbose(verbose bool) func(*cleaner) {
	return func(c *cleaner) {
		c.verbose = verbose
	}
}

func StdOut(out io.Writer) func(*cleaner) {
	return func(c *cleaner) {
		c.stdOut = out
	}
}
