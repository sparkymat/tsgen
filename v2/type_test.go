package v2_test

import "testing"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	IsStarred   bool    `json:"isStarred"`
	Author      User    `json:"author"`
	StarredTime *string `json:"starredTime"`
}

type PostCreateRequest struct {
	AuthorID string `json:"authorId"`
	Title    string `json:"title"`
}

type PostCreateResponse struct {
	Post   Post `json:"post"`
	Author User `json:"author"`
}

type PostSearchRequest struct {
	Query      string `query:"query"`
	PageSize   int32  `query:"pageSize"`
	PageNumber int32  `query:"pageNumber"`
}

type PostSearchResponse struct {
	Posts      []Post `json:"posts"`
	TotalCount int64  `json:"totalCount"`
}

func TestSlice(t *testing.T) {
	t.Run("complex type", func(t *testing.T) {
	})
}
