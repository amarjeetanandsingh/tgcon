package gen

import (
	"bytes"
	"strings"
	"testing"

	"github.com/amarjeetanandsingh/tgconst/parser"
)

// TODO:: to test it end to end: put testdata/data.go file in gen pkg.
var parsedFiles = []parser.File{
	{
		PackageName: "pkgName",
		FileName:    "fileName",
		Structs: []parser.Struct{
			{ // magic commented struct
				Name: "MagicCommentStruct",
				Fields: []parser.Field{
					{Name: "UntaggedField"},
					{Name: "OneTagField", Tags: []parser.Tag{{"tagKey", "tagVal"}}},
					{Name: "TwoTagField", Tags: []parser.Tag{{"tag1Key", "tag1Val"}, {"tag2Key", "tag2Val"}}},
				},
			},
		},
	},
}

func TestGenerateAndWrite(t *testing.T) {

	parsedFile := parser.File{}
	parsedFile.PackageName = "pkgName"
	parsedFile.FileName = "FileName.go"
	parsedFile.Structs = []parser.Struct{
		{ // magic commented struct
			Name: "MagicCommentStruct",
			Fields: []parser.Field{
				{Name: "UntaggedField"},
				{Name: "OneTagField", Tags: []parser.Tag{{"tagKey", "tagVal"}}},
				{Name: "TwoTagField", Tags: []parser.Tag{{"tag1Key", "tag1Val"}, {"tag2Key", "tag2Val"}}},
			},
		},
	}

	writer := &bytes.Buffer{}
	parsedFiles := []parser.File{parsedFile}
	g := New()
	if err := g.generateAndWrite(parsedFiles, writer); err != nil {
		t.Errorf("error writing parsed value to writer:: %w", err)
	}

	//

}

func TestTaggedFieldOnly(t *testing.T) {
	g := New()
	g.onlyTaggedFields = true

	writer := &bytes.Buffer{}
	if err := g.generateAndWrite(parsedFiles, writer); err != nil {
		t.Errorf("error in generateAndWrite. Error: %w", err)
	}
	output := writer.String()
	if strings.Contains(output, "UntaggedField") {
		t.Errorf("UntaggedField not expected")
	}
}

func TestAllowedTags(t *testing.T) {
	g := New()
	g.tags = []string{"tagKey"}

	writer := &bytes.Buffer{}
	if err := g.generateAndWrite(parsedFiles, writer); err != nil {
		t.Errorf("error in generateAndWrite. Error: %w", err)
	}
	output := writer.String()
	if strings.Contains(output, "tag1Key") {
		t.Errorf("tag1Key not expected")
	}

	if !strings.Contains(output, g.tags[0]) {
		t.Error(g.tags[0], " was not found")
	}
}

func TestName(t *testing.T) {

}
