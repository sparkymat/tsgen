package tsgen_test

import (
	"testing"

	"github.com/sparkymat/tsgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestSliceExport(t *testing.T) {
	t.Parallel()

	t.Run("simple slice", func(t *testing.T) {
		t.Parallel()

		s := tsgen.New()

		err := s.AddSliceEntry(
			"Post",
			"",
			"create",
			tsgen.ActionCreate,
			PostCreateRequest{},
			PostCreateResponse{},
			false,
			"",
		)
		require.NoError(t, err)

		contentMap, err := s.Export()
		require.NoError(t, err)

		assert.Equal(t, `import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

import { Post } from '../models/Post';
import { User } from '../models/User';

export interface PostCreateRequest {
  authorId: string;
  title: string;
}

export interface PostCreateResponse {
  post: Post;
  author: User;
}

export const api = createApi({
  reducerPath: 'posts',
  baseQuery: fetchBaseQuery({ baseUrl: '/api' }),
  tagTypes: ['Post'],
  endpoints: builder => ({
    create: builder.mutation<PostCreateResponse, PostCreateRequest>({
      query: request => ({
        url: `+"`"+`posts`+"`"+`,
        method: 'POST',
        body: request,
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content
        },
      }),
      invalidatesTags: [
        { type: 'Post', id: 'LIST' },
      ],
    }),
  })
});

export const {
  useCreateMutation,
} = api;
`, string(contentMap["slices/Post.ts"]))
	})

	t.Run("complex slice", func(t *testing.T) {
		t.Parallel()

		s := tsgen.New()

		err := s.AddModel(Post{})
		require.NoError(t, err)

		err = s.AddSliceEntry(
			"Post",
			"",
			"create",
			tsgen.ActionCreate,
			PostCreateRequest{},
			Post{},
			false,
			"",
		)
		require.NoError(t, err)

		err = s.AddSliceEntry(
			"Post",
			"",
			"search",
			tsgen.ActionCustomQuery,
			PostSearchRequest{},
			PostSearchResponse{},
			false,
			"",
		)
		require.NoError(t, err)

		contentMap, err := s.Export()
		require.NoError(t, err)

		assert.Equal(t, `import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

import { Post } from '../models/Post';
import { User } from '../models/User';

export interface PostSearchRequest {
  query: string;
  pageSize: number;
  pageNumber: number;
}

export interface PostSearchResponse {
	posts: Post[];
	totalCount: number;
}

export interface PostCreateResponse {
  post: Post;
  author: User;
}

export interface PostCreateRequest {
  authorId: string;
  title: string;
}

export const api = createApi({
  reducerPath: 'posts',
  baseQuery: fetchBaseQuery({ baseUrl: '/api' }),
  tagTypes: ['Post'],
  endpoints: builder => ({
    create: builder.mutation<PostCreateResponse, PostCreateRequest>({
      query: request => ({
        url: `+"`"+`posts`+"`"+`,
        method: 'POST',
        body: request,
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content
        },
      }),
      invalidatesTags: [
        { type: 'Post', id: 'LIST' },
      ],
    }),
  })
});

export const {
  useCreateMutation,
} = api;
`, string(contentMap["slices/Post.ts"]))
	})
}
