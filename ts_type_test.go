package tsgen_test

import (
	"testing"

	"github.com/sparkymat/tsgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrTo[T any](v T) *T { return &v }

func TestStructToTSType(t *testing.T) {
	t.Parallel()

	type User struct {
		ID         string  `json:"id"`
		Name       string  `json:"name"`
		PostCount  int     `json:"postCount"`
		Verfied    bool    `json:"verified"`
		UpgradedAt *string `json:"upgradedAt"`
	}

	type Comment struct {
		ID       string `json:"id"`
		Author   User   `json:"author"`
		Content  string `json:"content"`
		PostedAt string `json:"postedAt"`
	}

	type Post struct {
		ID       string    `json:"id"`
		Title    string    `json:"title"`
		Author   User      `json:"author"`
		Comments []Comment `json:"comments"`
	}

	tests := []struct {
		name           string
		v              any
		addID          bool
		expectedName   string
		expectedFields map[string]string
		wantErr        bool
	}{
		{
			name: "successfully convert valid struct",
			v: func() any {
				val := User{
					ID:         "1",
					Name:       "Jack Frost",
					PostCount:  175,
					Verfied:    false,
					UpgradedAt: ptrTo("2021-11-16T05:00:00Z"),
				}

				return val
			}(),
			expectedName: "User",
			expectedFields: map[string]string{
				"id":          "string",
				"name":        "string",
				"postCount":   "number",
				"upgradedAt?": "string",
				"verified":    "boolean",
			},
			wantErr: false,
		},
		{
			name: "successfully convert valid struct and inject id",
			v: func() any {
				val := User{
					Name:       "Jack Frost",
					PostCount:  175,
					Verfied:    false,
					UpgradedAt: ptrTo("2021-11-16T05:00:00Z"),
				}

				return val
			}(),
			addID:        true,
			expectedName: "User",
			expectedFields: map[string]string{
				"id":          "string",
				"name":        "string",
				"postCount":   "number",
				"upgradedAt?": "string",
				"verified":    "boolean",
			},
			wantErr: false,
		},
		{
			name: "successfully convert valid nested struct",
			v: func() any {
				val := Post{
					ID:       "5",
					Title:    "Converting to TS",
					Author:   User{},
					Comments: []Comment{},
				}

				return val
			}(),
			expectedName: "Post",
			expectedFields: map[string]string{
				"id":       "string",
				"title":    "string",
				"author":   "User",
				"comments": "Comment[]",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output, err := tsgen.StructToTSType(tt.v, tt.addID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedName, output.Name())
				assert.Equal(t, tt.expectedFields, output.Fields())
			}
		})
	}
}
