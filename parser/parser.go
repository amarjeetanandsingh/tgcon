package parser

import (
	"fmt"
	"go/ast"
	goParser "go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type parser struct {

	/*
		`tgconst` is a magic comment which makes a struct eligible to be parsed.
		In the code below, fields of struct Foo, Bar, Baz and Qux will be parsed because
		they have magic comment(tgconst) associated with them.

		```
		// tgconst
		type Foo struct {
			F1 int `json:"f1"`
		}

		// tgconst
		type (
			Bar struct {
				F1 int `json:"f1"`
			}
			Baz struct {
				F1 int `json:"f1"`
			}
		)
		type (
			// tgconst
			Qux struct {
				F1 int `json:"f1"`
			}
			Quuz struct {
				F1 int `json:"f1"`
			}
		)
		type Xyz struct {
			F1 int `json:"f1"`
		}
		```

		Struct Quuz and Xyz are not parsed because they don't have the magic comment.
	*/
	magicComment string

	// len(tags) == 0 means create const for each tags associated to a field.
	// Set of allowed tags to be processed.
	tags map[string]bool

	// consider only tagged fields
	onlyTaggedFields bool

	// process all structs irrespective of magic comment
	allStructs bool
}

func New(magicComment string, tags []string, allStructs, onlyTaggedFields bool) *parser {
	p := &parser{
		allStructs:       allStructs,
		magicComment:     magicComment,
		onlyTaggedFields: onlyTaggedFields,
		tags:             map[string]bool{},
	}

	// allowed tags
	for _, t := range tags {
		p.tags[t] = true
	}
	return p
}

func (p parser) ParseDir(dir string) ([]File, error) {
	filesInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory: %s :: %w", dir, err)
	}

	var parsedFiles []File

	for _, file := range filesInfo {
		// skip invalid files
		if file.IsDir() || // skip directory
			len(file.Name()) == 0 ||
			file.Name()[0] == '.' || // hidden file
			!strings.HasSuffix(file.Name(), ".go") || // not a go file
			strings.HasSuffix(file.Name(), "_gen.go") || // not a generated file
			strings.HasSuffix(file.Name(), "_test.go") { // not a test file
			continue
		}

		// parse file
		filePath := path.Join(dir, file.Name())
		reader, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %s :: %w", filePath, err)
		}

		parsedFile, err := p.ParseFile(reader)
		reader.Close()
		if err != nil {
			return nil, fmt.Errorf("error parsing file: %s :: %w", filePath, err)
		}

		if parsedFile != nil {
			parsedFile.FileName = file.Name()
			parsedFiles = append(parsedFiles, *parsedFile)
		}
	}

	return parsedFiles, nil
}

// ParseFile accepts a reader to a go source code and returns
// nil, nil if the source doesn't have a valid parsed struct.
func (p parser) ParseFile(reader io.Reader) (parsedFile *File, err error) {
	astFile, err := goParser.ParseFile(token.NewFileSet(), "", reader, goParser.ParseComments)
	if err != nil {
		return nil, err
	}

	parsedFile = &File{}
	parsedFile.PackageName = astFile.Name.Name
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.GenDecl:
			commentOnTypeDec := strings.TrimSpace(v.Doc.Text())
			isTypeDecAnnotated := strings.Contains(commentOnTypeDec, p.magicComment)

			for _, spec := range v.Specs {

				typeSpec, isTypeSpec := spec.(*ast.TypeSpec)
				if !isTypeSpec {
					continue
				}

				// !p.allStructs:- process all structs even if magic comment is not there.
				// !isTypeDecAnnotated:- typeDec doesn't have the magic comment
				// !strings.Contains(typeSpec.Doc.Text(), p.magicComment):- struct doesn't have magic comment
				if !p.allStructs && !isTypeDecAnnotated && !strings.Contains(typeSpec.Doc.Text(), p.magicComment) {
					continue
				}

				parsedStruct, e := p.parseStruct(typeSpec)
				if e != nil {
					err = e
					return false
				}
				if parsedStruct != nil && len(parsedStruct.Fields) > 0 {
					parsedFile.Structs = append(parsedFile.Structs, *parsedStruct)
				}
			}
		}
		return true
	})

	// if no struct, do not return file
	if parsedFile == nil || len(parsedFile.Structs) == 0 {
		return nil, err
	}

	// sort structs by name
	sort.Slice(parsedFile.Structs, func(i, j int) bool {
		return parsedFile.Structs[i].Name < parsedFile.Structs[j].Name
	})
	return parsedFile, err
}

func (p parser) parseStruct(typeSpec *ast.TypeSpec) (*Struct, error) {
	str, isStructType := typeSpec.Type.(*ast.StructType)
	if !isStructType {
		return nil, nil
	}

	parsedStruct := &Struct{}
	parsedStruct.Name = typeSpec.Name.Name

	for _, field := range str.Fields.List {
		var parsedTags []Tag
		if field.Tag != nil {
			parsedTags = p.parseTag(field.Tag.Value)
		}

		// Skip if non tagged fields are not allowed
		// If multiple fields are there a in single line, ignore tags.
		if p.onlyTaggedFields && (len(parsedTags) == 0 || len(field.Names) > 1) {
			continue
		}

		// type strName struct {
		// 		Field1, Field2, Field3 string `json:"json,omitempty"`
		// }
		// Ignore tag in this case because we can't decide
		// on which field json tag is attached to.
		if len(field.Names) > 1 {
			for _, fieldName := range field.Names {
				parsedField := Field{}
				parsedField.Name = fieldName.Name
				parsedStruct.Fields = append(parsedStruct.Fields, parsedField)
			}
			continue
		}

		// add fields even if there are no tags.
		parsedStruct.Fields = append(parsedStruct.Fields, Field{
			Name: field.Names[0].Name,
			Tags: parsedTags,
		})
	}

	// do not return struct if not fields are there.
	if len(parsedStruct.Fields) == 0 {
		return nil, nil
	}

	// sort fields
	sort.Slice(parsedStruct.Fields, func(i, j int) bool {
		return parsedStruct.Fields[i].Name < parsedStruct.Fields[j].Name
	})

	return parsedStruct, nil
}

func (p parser) parseTag(tag string) []Tag {
	var tags []Tag

	// Skip leading back-tick or double quotes coming from *ast.field.Tag.Value literal.
	// field.Tag.Value gave tag value as `json: "val"`
	// This use to return `json as key.
	for tag != "" && (tag[0] == '`' || tag[0] == '"') {
		tag = tag[1:]
	}

	n := 0
	for ; tag != ""; n++ {
		if n > 0 && tag != "" && tag[0] != ' ' {
			// More restrictive than reflect, but catches likely mistakes
			// like `x:"foo",y:"bar"`, which parses as `x:"foo" ,y:"bar"` with second key ",y".
			break
		}
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 {
			break
		}
		if i+1 >= len(tag) || tag[i] != ':' {
			break
		}
		if tag[i+1] != '"' {
			break
		}
		key := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		value = strings.Split(value, ",")[0]

		// len(p.tags) == 0 means all tags allowed
		// p.tags[key] == true means tag named `key` is allowed
		// Only add a tag if there is a value. Treat it as untagged otherwise.
		if (len(p.tags) == 0 || p.tags[key]) && value != "" {
			tags = append(tags, Tag{
				Name:  key,
				Value: value,
			})
		}
	}

	// sort tags
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	return tags
}
