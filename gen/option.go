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

import "github.com/amarjeetanandsingh/tgconst/text"

// moved the option functions from gen.go to keep the logic clean out there

func All(all bool) func(*generator) {
	return func(g *generator) {
		g.allStructs = all
	}
}

func Dir(d string) func(*generator) {
	return func(g *generator) {
		g.dir = d
	}
}

func Tags(t []string) func(*generator) {
	return func(g *generator) {
		g.tags = t
	}
}

func Recursive(r bool) func(*generator) {
	return func(g *generator) {
		g.isRecursive = r
	}
}

func TaggedFieldOnly(tagged bool) func(*generator) {
	return func(g *generator) {
		g.onlyTaggedFields = tagged
	}
}

func MissingTagPolicy(policy string) func(*generator) {
	return func(g *generator) {
		switch policy {
		case text.SnakeCase:
			g.missingTagValFormat = text.SnakeCase
		case text.CamelCase:
			g.missingTagValFormat = text.CamelCase
		case text.LispCase:
			g.missingTagValFormat = text.LispCase
		case text.PascalCase:
			g.missingTagValFormat = text.PascalCase
		default:
			g.missingTagValFormat = text.Mirror
		}
	}
}
