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
		transFormScheme TransformScheme
	}{
		{input: "", want: ""},
		{input: "ONETWO", want: "onetwo", transFormScheme: SnakeCase},
		{input: "OneTwo", want: "one_two", transFormScheme: SnakeCase},
		{input: "oneTwo", want: "one_two", transFormScheme: SnakeCase},
		{input: "onetwo", want: "onetwo", transFormScheme: SnakeCase},
		{input: "ONE2", want: "one_2", transFormScheme: SnakeCase},
		{input: "One2", want: "one_2", transFormScheme: SnakeCase},
		{input: "one2", want: "one_2", transFormScheme: SnakeCase},

		{input: "ONETWO", want: "onetwo", transFormScheme: CamelCase},
		{input: "OneTwo", want: "oneTwo", transFormScheme: CamelCase},
		{input: "oneTwo", want: "oneTwo", transFormScheme: CamelCase},
		{input: "onetwo", want: "onetwo", transFormScheme: CamelCase},
		{input: "ONE2", want: "one2", transFormScheme: CamelCase},
		{input: "One2", want: "one2", transFormScheme: CamelCase},
		{input: "one2", want: "one2", transFormScheme: CamelCase},

		{input: "ONETWO", want: "ONETWO", transFormScheme: PascalCase},
		{input: "OneTwo", want: "OneTwo", transFormScheme: PascalCase},
		{input: "oneTwo", want: "OneTwo", transFormScheme: PascalCase},
		{input: "onetwo", want: "Onetwo", transFormScheme: PascalCase},
		{input: "ONE2", want: "ONE2", transFormScheme: PascalCase},
		{input: "One2", want: "One2", transFormScheme: PascalCase},
		{input: "one2", want: "One2", transFormScheme: PascalCase},
	} {
		if got := Transform(test.input, test.transFormScheme); got != test.want {
			t.Errorf("got: %s, want: %s, scheme: %s", got, test.want, test.transFormScheme)
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
