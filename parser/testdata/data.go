package testdata

// tgconst
type MagicCommentedStruct1 struct {
	UnTaggedField    int
	SingleTagField   int `json:"signeTagVal"`
	MultiTaggedField int `json:"jsonTagVal" fooTag:"fooTagVal"`
}

// tgconst
type (
	Bar struct {
		F1 int `json:"bartag"`
	}
	Baz struct {
		F1 int `json:"baztag"`
	}
)

type (
	// tgconst
	Qux struct {
		F1 int `json:"quxtag"`
	}
	Quuz struct {
		F1 int `json:"quuztag"`
	}
)

type Xyz struct {
	F1 int `json:"xyz_tag"`
}

// tgconst
type StructA struct {
	NameA string `json:"name_a" xml:"name_a_xml"`
	AgeA  int    `json:"age_a" xml:"age_a_xml"`
}

type StructB struct {
	NameB string
	AgeB  int
}

// tgconst
type (
	// StructC struct comment
	StructC struct {
		NameC string `json:"name_c"`
		AgeC  int
	}
)
