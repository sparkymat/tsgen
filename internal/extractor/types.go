package extractor

// StructType contains the details of a struct type
type StructType struct {
	// Name of the struct
	Name string
	// Fields information
	Fields []StructField
}

// StructField contains the details of a struct field
type StructField struct {
	// Name of the field
	Name string
	// Tags found on the field (tag split by space)
	Tags []string
	// TagMap contains the tags extracted and mapped by type (e.g. json:"foo,required"
	// will become map[string][]string{"json", {"foo","required"}}
	TagMap map[string][]string
	// Type of the field
	Type Type
	// TypeName will contain the fully qualified type for non-basic types
	TypeName string
}

// Type is the type of a field
type Type string

const (
	TypeString  Type = "string"
	TypeNumber  Type = "number"
	TypeBoolean Type = "boolean"
	TypeStruct  Type = "struct"
	TypeUnknown Type = "unknown"
)
