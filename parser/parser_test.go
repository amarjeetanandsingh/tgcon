package parser

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestParseFile(t *testing.T) {
	testFilePath := path.Join("testdata", "data.go")
	reader, err := os.Open(testFilePath)
	if err != nil {
		t.Errorf("error reading test file: %s :: %w", testFilePath, err)
	}

	p := New("tgconst", nil, false, true)
	file, err := p.ParseFile(reader)
	if err != nil {
		t.Errorf("error parsing file: %w", err)
	}
	log.Println(file)
}
