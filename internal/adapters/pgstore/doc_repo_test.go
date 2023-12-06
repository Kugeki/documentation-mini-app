package pgstore

import (
	"context"
	"documentation-mini-app/internal/domain/article"
	"documentation-mini-app/internal/domain/doc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocRepoPG_CreateDoc(t *testing.T) {
	ctx := context.TODO()

	s, truncate := TestStore(ctx, t, dbURL)
	defer truncate(ctx, "documentation")

	d := doc.Documentation{
		Name:                     "example",
		DefaultHighlightLanguage: "",
	}

	err := s.Doc().Create(ctx, &d)

	assert.NoError(t, err)
	assert.NotZero(t, d.ID)
}

func TestDocRepoPG_GetDoc(t *testing.T) {
	ctx := context.TODO()

	s, truncate := TestStore(ctx, t, dbURL)
	defer truncate(ctx, "documentation")

	docID := 1

	_, err := s.Doc().GetByID(ctx, docID)
	assert.Error(t, err)

	d := doc.Documentation{
		Name:                     "example",
		DefaultHighlightLanguage: "Go",
		Articles:                 []article.Article{},
	}

	err = s.Doc().Create(ctx, &d)
	assert.NoError(t, err)

	getD, err := s.Doc().GetByID(ctx, docID)
	assert.NoError(t, err)
	assert.Equal(t, *getD, d)
}

func TestDocRepoPG_Update(t *testing.T) {
	ctx := context.TODO()

	s, truncate := TestStore(ctx, t, dbURL)
	defer truncate(ctx, "documentation")

	d := doc.Documentation{
		ID:                       1,
		Name:                     "example",
		DefaultHighlightLanguage: "Go",
		Articles:                 []article.Article{},
	}

	err := s.Doc().Update(ctx, &d)
	assert.Error(t, err)

	err = s.Doc().Create(ctx, &d)
	assert.NoError(t, err)

	newName := "updated_example"
	d.Name = newName

	err = s.Doc().Update(ctx, &d)
	assert.NoError(t, err)

	getD, err := s.Doc().GetByID(ctx, d.ID)
	assert.NoError(t, err)
	assert.Equal(t, *getD, d)
}

func TestDocRepoPG_DeleteDoc(t *testing.T) {
	ctx := context.TODO()

	s, truncate := TestStore(ctx, t, dbURL)
	defer truncate(ctx, "documentation")

	docID := 1

	err := s.Doc().Delete(ctx, docID)
	assert.Error(t, err)

	d := doc.Documentation{
		Name:                     "example",
		DefaultHighlightLanguage: "",
	}

	err = s.Doc().Create(ctx, &d)
	assert.NoError(t, err)

	err = s.Doc().Delete(ctx, d.ID)
	assert.NoError(t, err)

	getDoc, err := s.Doc().GetByID(ctx, d.ID)
	assert.Error(t, err)
	assert.Nil(t, getDoc)
}
