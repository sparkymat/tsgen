package tsgen_test

import (
	"testing"

	"github.com/sparkymat/tsgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Foo struct{}

type User struct {
	ID string `json:"id"`
}

type Post struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	IsStarred   bool    `json:"isStarred"`
	Author      User    `json:"author"`
	StarredTime *string `json:"starredTime"`
}

func TestModelExport(t *testing.T) {
	t.Parallel()

	t.Run("empty type", func(t *testing.T) {
		t.Parallel()

		s := tsgen.New()

		err := s.AddModel(Foo{})
		require.NoError(t, err)

		contentMap, err := s.Export()
		require.NoError(t, err)

		assert.Equal(t, `export class Foo {
  constructor(json: Foo) {
  }
}
`, string(contentMap["models/Foo.ts"]))
	})

	t.Run("valid type", func(t *testing.T) {
		t.Parallel()

		s := tsgen.New()

		err := s.AddModel(User{})
		require.NoError(t, err)

		contentMap, err := s.Export()
		require.NoError(t, err)

		assert.Equal(t, `export class User {
  public id: string;

  constructor(json: User) {
    this.id = json.id;
  }
}
`, string(contentMap["models/User.ts"]))
	})

	t.Run("valid type with non-basic types", func(t *testing.T) {
		t.Parallel()

		s := tsgen.New()

		err := s.AddModel(Post{})
		require.NoError(t, err)

		contentMap, err := s.Export()
		require.NoError(t, err)

		assert.Equal(t, `import { User } from './User';

export class Post {
  public id: string;

  public title: string;

  public isStarred: boolean;

  public author: User;

  public starredTime?: string;

  constructor(json: Post) {
    this.id = json.id;
    this.title = json.title;
    this.isStarred = json.isStarred;
    this.author = new User(json.author);
    if (json.starredTime) {
      this.starredTime = json.starredTime;
    }
  }
}
`, string(contentMap["models/Post.ts"]))
	})
}
