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

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path"

	"github.com/amarjeetanandsingh/tgcon/files"
	"github.com/amarjeetanandsingh/tgcon/parser"
	"github.com/amarjeetanandsingh/tgcon/text"
)

type generator struct {
	// List of tags we are going to create constants for
	tags []string

	// Do not create string constants for those fields which doesn't have tags.
	// Ex: No const will be generated for F1 field of `struct{F1 int}` because
	// it has no tag.
	onlyTaggedFields bool // TODO:v1: defaultTagVal can be better name

	// todo: v2 :
	// 	By default it adds tag name as suffix to the const variable. When we have
	//	more than one tags to a field, its required to append the tag name to the
	//	end of the const variable to create two different tags.
	//	Eg:
	//		type Str struct{ F1 int `json:"f1" bson:"f1"`}
	//	Will return two const, 1) Str_F1_json, 2) Str_F1_bson
	//	If noSuffix = true, it will create only one const per struct field, as `Str_F1`
	noSuffix bool

	// Setting `isRecursive` true will generate const for all the subdirectories too.
	// By default it generates const only for the current directory.
	isRecursive bool

	// Generate const file for the given directory. If `isRecursive` flag is set,
	// it will generate const file recursively for all its subdirectories too.
	dir string

	// process all structs irrespective of magic comment
	allStructs bool

	// todo: better name and doc
	// Format to generate tag value for untagged fields
	// possible values are [CamelCase, LispCase, PascalCase, SnakeCase, Mirror]
	// Mirror is default value in case no(empty) TransformFormat was given
	defaultTagValFormat text.TransformFormat
}

func New(options ...func(*generator)) *generator {
	g := &generator{}

	// set config
	for _, option := range options {
		option(g)
	}

	// validate options
	if g.onlyTaggedFields && len(g.defaultTagValFormat) > 0 {
		fmt.Println("Warning: defaultTagValFormat ignored because onlyTagged flag is set")
	}

	return g
}

func (g *generator) Do() error {
	return g.generateConstantsFile(g.dir)
}

func (g *generator) generateConstantsFile(dir string) error {
	p := parser.New("tgcon", g.tags, g.allStructs, g.onlyTaggedFields)
	parsedFiles, err := p.ParseDir(dir)
	if err != nil {
		return err
	}

	// only write to a file if there is some struct field
	if len(parsedFiles) > 0 {
		generatedFilePath := path.Join(dir, parsedFiles[0].PackageName+"_tgcon_gen.go")
		generatedFileWriter, err := os.OpenFile(generatedFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("error creating generated file %s: "+err.Error(), generatedFilePath)
		}
		if err := g.generateAndWrite(parsedFiles, generatedFileWriter); err != nil {
			return fmt.Errorf("error writing to generated file %s :: "+err.Error(), generatedFilePath)
		}
	}

	// recursive check
	if !g.isRecursive {
		return nil
	}

	// if is recursive, continue with sub-directories
	subDirNames, err := files.ListDirs(dir, func(dirName string) bool {
		// TODO comment this in doc
		return len(dirName) > 0 && dirName[0] != '.'
	})
	if err != nil {
		return fmt.Errorf("error listing subdirs of %s :: "+err.Error(), dir)
	}
	for _, subDir := range subDirNames {
		if err := g.generateConstantsFile(path.Join(dir, subDir)); err != nil {
			return err
		}
	}
	return nil
}

// todo better name?
func (g *generator) generateAndWrite(parsedFiles []parser.File, writer io.Writer) error {
	if len(parsedFiles) == 0 {
		return nil
	}

	formattedCode, err := g.generateFormattedCode(parsedFiles)
	if err != nil {
		return err
	}

	bw := bufio.NewWriter(writer)
	if _, err := bw.Write(formattedCode); err != nil {
		return fmt.Errorf("error writing formatted code to writer :: " + err.Error())
	}
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("error flushing formatted code to writer :: " + err.Error())
	}

	return nil
}

func (g *generator) generateFormattedCode(parsedFiles []parser.File) ([]byte, error) {
	generatedCode, err := g.generateCode(parsedFiles)
	if err != nil {
		return nil, err
	}
	// TODO VERIFY IF THIS CHECK IS NEEDED
	if len(generatedCode) == 0 {
		return nil, nil
	}

	formattedSource, err := format.Source(generatedCode)
	if err != nil {
		return nil, fmt.Errorf("error formatting generated source code :: " + err.Error())
	}

	return formattedSource, nil
}

func (g *generator) generateCode(parsedFiles []parser.File) ([]byte, error) {
	if len(parsedFiles) == 0 {
		return nil, nil
	}

	sourceCode := bytes.Buffer{}
	sourceCode.WriteString("// Code generated by tgcon; DO NOT EDIT.\n")
	sourceCode.WriteString(fmt.Sprintf("package %s\n\n", parsedFiles[0].PackageName))
	sourceCode.WriteString("const (\n\n")

	noConstants := true
	for _, parsedFile := range parsedFiles {
		sourceCode.WriteString("// " + text.CenterAlignedPadded("File: "+parsedFile.FileName, "-") + "\n")

		for _, parsedStruct := range parsedFile.Structs {
			sourceCode.WriteString("// Struct: " + parsedStruct.Name + "\n")

			for _, parsedField := range parsedStruct.Fields {
				for _, tag := range parsedField.Tags {
					noConstants = false // means at least 1 const is added.

					line := g.generateLine(parsedStruct.Name, parsedField.Name, tag.Name, tag.Value)
					sourceCode.WriteString(line + "\n")
				}

				// means generate tag constant value for parsedField which has no tags
				if len(parsedField.Tags) == 0 && !g.onlyTaggedFields {
					noConstants = false

					line := g.generateLine(parsedStruct.Name, parsedField.Name, "", parsedField.Name)
					sourceCode.WriteString(line + "\n")
				}
			}
			sourceCode.WriteString("\n")
		}
		sourceCode.WriteString("\n")
	}
	sourceCode.WriteString(")\n")

	// if no field constant is written, do not return any file content.
	if noConstants {
		return nil, nil
	}

	return sourceCode.Bytes(), nil
}

func (g *generator) generateLine(structName, FieldName, tagName, tagValue string) string {
	words := make([]string, 0, 3)

	if structName != "" {
		words = append(words, structName)
	}
	if FieldName != "" {
		words = append(words, FieldName)
	}
	if tagName != "" {
		words = append(words, tagName)
	}

	constVariable := text.JoinFormatted(words, text.PascalCase) // TODO:vote: take from g.constVariableFormat
	constValue := text.Format(tagValue, g.defaultTagValFormat)
	return constVariable + " = \"" + constValue + "\""
}
