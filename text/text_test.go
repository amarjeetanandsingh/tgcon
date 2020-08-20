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
		transformFormat TransformFormat
	}{
		{input: "", want: ""},
		{input: "ONETWO", want: "onetwo", transformFormat: SnakeCase},
		{input: "OneTwo", want: "one_two", transformFormat: SnakeCase},
		{input: "oneTwo", want: "one_two", transformFormat: SnakeCase},
		{input: "onetwo", want: "onetwo", transformFormat: SnakeCase},
		{input: "ONE2", want: "one_2", transformFormat: SnakeCase},
		{input: "One2", want: "one_2", transformFormat: SnakeCase},
		{input: "one2", want: "one_2", transformFormat: SnakeCase},

		{input: "ONETWO", want: "onetwo", transformFormat: CamelCase},
		{input: "OneTwo", want: "oneTwo", transformFormat: CamelCase},
		{input: "oneTwo", want: "oneTwo", transformFormat: CamelCase},
		{input: "onetwo", want: "onetwo", transformFormat: CamelCase},
		{input: "ONE2", want: "one2", transformFormat: CamelCase},
		{input: "One2", want: "one2", transformFormat: CamelCase},
		{input: "one2", want: "one2", transformFormat: CamelCase},

		{input: "ONETWO", want: "ONETWO", transformFormat: PascalCase},
		{input: "OneTwo", want: "OneTwo", transformFormat: PascalCase},
		{input: "oneTwo", want: "OneTwo", transformFormat: PascalCase},
		{input: "onetwo", want: "Onetwo", transformFormat: PascalCase},
		{input: "ONE2", want: "ONE2", transformFormat: PascalCase},
		{input: "One2", want: "One2", transformFormat: PascalCase},
		{input: "one2", want: "One2", transformFormat: PascalCase},

		{input: "ONETWO", want: "onetwo", transformFormat: LispCase},
		{input: "OneTwo", want: "one-two", transformFormat: LispCase},
		{input: "oneTwo", want: "one-two", transformFormat: LispCase},
		{input: "onetwo", want: "onetwo", transformFormat: LispCase},
		{input: "ONE2", want: "one-2", transformFormat: LispCase},
		{input: "One2", want: "one-2", transformFormat: LispCase},
		{input: "one2", want: "one-2", transformFormat: LispCase},
	} {
		if got := Format(test.input, test.transformFormat); got != test.want {
			t.Errorf("got: %s, want: %s, scheme: %s", got, test.want, test.transformFormat)
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
