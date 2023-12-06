package articleuc

import (
	"context"
	"documentation-mini-app/internal/adapters/pgstore"
	"documentation-mini-app/internal/domain/article"
)

type ArticleUC struct {
	Store *pgstore.Store
}

func New(store *pgstore.Store) *ArticleUC {
	return &ArticleUC{Store: store}
}

func (uc *ArticleUC) GetArticleByID(ctx context.Context, id int) (*article.Article, error) {
	art, err := uc.Store.Article().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return art, nil
}

func (uc *ArticleUC) CreateArticle(ctx context.Context, art *article.Article) error {
	err := uc.Store.Article().Create(ctx, art)
	return err
}

func (uc *ArticleUC) AddArticleToDoc(ctx context.Context, artID int, docID int) error {
	if docID != 0 {
		err := uc.Store.Article().AddToDoc(ctx, artID, docID)
		return err
	}
	return nil
}

func (uc *ArticleUC) UpdateArticle(ctx context.Context, art *article.Article) error {
	err := uc.Store.Article().Update(ctx, art)
	return err
}

func (uc *ArticleUC) DeleteArticle(ctx context.Context, artID int) error {
	err := uc.Store.Article().Delete(ctx, artID)
	return err
}
