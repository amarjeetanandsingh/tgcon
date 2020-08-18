package text

import "testing"

func TestSplit(t *testing.T) {
	for _, test := range []struct {
		input string
		want  []string
	}{
		{input: "", want: []string{}},
		{input: "lowercase", want: []string{"lowercase"}},
		{input: "Class", want: []string{"Class"}},
		{input: "MyClass", want: []string{"My", "Class"}},
		{input: "MyC", want: []string{"My", "C"}},
		{input: "HTML", want: []string{"HTML"}},
		{input: "PDFLoader", want: []string{"PDF", "Loader"}},
		{input: "AString", want: []string{"A", "String"}},
		{input: "SimpleXMLParser", want: []string{"Simple", "XML", "Parser"}},
		{input: "vimRPCPlugin", want: []string{"vim", "RPC", "Plugin"}},
		{input: "GL11Version", want: []string{"GL", "11", "Version"}},
		{input: "99Bottles", want: []string{"99", "Bottles"}},
		{input: "May5", want: []string{"May", "5"}},
		{input: "BFG9000", want: []string{"BFG", "9000"}},
		{input: "BöseÜberraschung", want: []string{"Böse", "Überraschung"}},
		{input: "Two  spaces", want: []string{"Two", "  ", "spaces"}},
		{input: "BadUTF8\xe2\xe2\xa1", want: []string{"BadUTF8\xe2\xe2\xa1"}},
	} {
		if got := Split(test.input); !equals(got, test.want) {
			t.Errorf("got: %v want: %v", got, test.want)
		}
	}
}

func TestTransform(t *testing.T) {
	for _, test := range []struct {
		input           string
		want            string
		transFormFoamat TransformFormat
	}{
		{input: "", want: ""},
		{input: "ONETWO", want: "onetwo", transFormFoamat: SnakeCase},
		{input: "OneTwo", want: "one_two", transFormFoamat: SnakeCase},
		{input: "oneTwo", want: "one_two", transFormFoamat: SnakeCase},
		{input: "onetwo", want: "onetwo", transFormFoamat: SnakeCase},
		{input: "ONE2", want: "one_2", transFormFoamat: SnakeCase},
		{input: "One2", want: "one_2", transFormFoamat: SnakeCase},
		{input: "one2", want: "one_2", transFormFoamat: SnakeCase},

		{input: "ONETWO", want: "onetwo", transFormFoamat: CamelCase},
		{input: "OneTwo", want: "oneTwo", transFormFoamat: CamelCase},
		{input: "oneTwo", want: "oneTwo", transFormFoamat: CamelCase},
		{input: "onetwo", want: "onetwo", transFormFoamat: CamelCase},
		{input: "ONE2", want: "one2", transFormFoamat: CamelCase},
		{input: "One2", want: "one2", transFormFoamat: CamelCase},
		{input: "one2", want: "one2", transFormFoamat: CamelCase},

		{input: "ONETWO", want: "ONETWO", transFormFoamat: PascalCase},
		{input: "OneTwo", want: "OneTwo", transFormFoamat: PascalCase},
		{input: "oneTwo", want: "OneTwo", transFormFoamat: PascalCase},
		{input: "onetwo", want: "Onetwo", transFormFoamat: PascalCase},
		{input: "ONE2", want: "ONE2", transFormFoamat: PascalCase},
		{input: "One2", want: "One2", transFormFoamat: PascalCase},
		{input: "one2", want: "One2", transFormFoamat: PascalCase},
	} {
		if got := Transform(test.input, test.transFormFoamat); got != test.want {
			t.Errorf("got: %s, want: %s, scheme: %s", got, test.want, test.transFormFoamat)
		}
	}
}

func equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
