package extractor_test

import (
	"testing"

	"github.com/sparkymat/tsgen/internal/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractStructTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		fileName      string
		fileContents  string
		errorExpected bool
		expectedTypes []extractor.StructType
	}{
		{
			name:     "extract 1 type with basic typed fields",
			fileName: "foobar.go",
			fileContents: `
package foo

type Response struct {
	Name string ` + "`" + `json:"name"` + "`" + `
	Age int64 ` + "`" + `json:"age"` + "`" + `
	IsRegistered bool ` + "`" + `json:"isRegistered"` + "`" + `
}
`,
			errorExpected: false,
			expectedTypes: []extractor.StructType{
				{
					Name: "Response",
					Fields: []extractor.StructField{
						{
							Name: "Name",
							Tags: []string{"json:\"name\""},
							TagMap: map[string][]string{
								"json": {"name"},
							},
							Type: extractor.TypeString,
						},
						{
							Name: "Age",
							Tags: []string{"json:\"age\""},
							TagMap: map[string][]string{
								"json": {"age"},
							},
							Type: extractor.TypeNumber,
						},
						{
							Name: "IsRegistered",
							Tags: []string{"json:\"isRegistered\""},
							TagMap: map[string][]string{
								"json": {"isRegistered"},
							},
							Type: extractor.TypeBoolean,
						},
					},
				},
			},
		},
		{
			name:     "extract 2 types with complex typed fields",
			fileName: "foobar.go",
			fileContents: `
package foo

type User struct {
	Name string ` + "`" + `json:"name"` + "`" + `
	Age int64 ` + "`" + `json:"age"` + "`" + `
	IsRegistered bool ` + "`" + `json:"isRegistered"` + "`" + `
}

type RelationEntry struct {
	User1 User ` + "`" + `json:"user1"` + "`" + `
	User2 User ` + "`" + `json:"user2"` + "`" + `
	RelationType RelationType ` + "`" + `json:"relationType"` + "`" + `
}
`,
			errorExpected: false,
			expectedTypes: []extractor.StructType{
				{
					Name: "User",
					Fields: []extractor.StructField{
						{
							Name: "Name",
							Tags: []string{"json:\"name\""},
							TagMap: map[string][]string{
								"json": {"name"},
							},
							Type: extractor.TypeString,
						},
						{
							Name: "Age",
							Tags: []string{"json:\"age\""},
							TagMap: map[string][]string{
								"json": {"age"},
							},
							Type: extractor.TypeNumber,
						},
						{
							Name: "IsRegistered",
							Tags: []string{"json:\"isRegistered\""},
							TagMap: map[string][]string{
								"json": {"isRegistered"},
							},
							Type: extractor.TypeBoolean,
						},
					},
				},
				{
					Name: "RelationEntry",
					Fields: []extractor.StructField{
						{
							Name: "User1",
							Tags: []string{"json:\"user1\""},
							TagMap: map[string][]string{
								"json": {"user1"},
							},
							Type: extractor.TypeStruct,
						},
						{
							Name: "User2",
							Tags: []string{"json:\"user2\""},
							TagMap: map[string][]string{
								"json": {"user2"},
							},
							Type: extractor.TypeStruct,
						},
						{
							Name: "RelationType",
							Tags: []string{"json:\"relationType\""},
							TagMap: map[string][]string{
								"json": {"relationType"},
							},
							Type: extractor.TypeBoolean,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			types, err := extractor.ExtractStructTypes(test.fileName, test.fileContents)

			if test.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedTypes, types)
			}
		})
	}
}
