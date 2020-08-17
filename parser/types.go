package parser

// Represents a go file. It contains tag info of all the eligible struct fields.
type File struct {
	FileName    string
	PackageName string
	Structs     []Struct
}

// Represents a single struct which had magic comment on it.
type Struct struct {
	Name   string
	Fields []Field
}

// Represents a singlle field of a struct
type Field struct {
	Name string
	Tags []Tag
}

// Represents a single tag of a struct field.
type Tag struct {
	Name  string
	Value string
}
