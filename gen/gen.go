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

package gen

import "fmt"

type generator struct {
	// List of tags we are going to create constants for
	tags []string

	// Do not create string constants for those fields which doesn't have tags.
	// Ex: No const will be generated for F1 field of `struct{F1 int}` because
	// it has no tag.
	onlyTaggedFields bool

	// Do not add tag name as suffix to the generated constant.
	// By default it will _tagName to the const.
	//
	// Ex:
	//		type Str struct{ F1 int `json:"f1"`}
	// Generated const will be `const Str_F1_json string = "f1"`
	noSuffix bool

	// Setting `isRecursive` true will generate const for all the subdirectories too.
	// By default it generates const only for the current directory.
	isRecursive bool

	// Generate const file for the given directory. If `isRecursive` flag is set,
	// it will generate const file recursively for all its subdirectories too.
	dir string
}

func New(options ...func(*generator)) *generator {
	g := &generator{}

	// set config
	for _, option := range options {
		option(g)
	}

	return g
}

func (g *generator) Do() error {
	fmt.Printf("generator: %+v", g)
	return nil
}
