package appuc

import (
	"context"
	"documentation-mini-app/internal/adapters/pgstore"
	"documentation-mini-app/internal/domain/article"
	"documentation-mini-app/internal/domain/crossed"
	"documentation-mini-app/internal/domain/doc"
)

type AppUC struct {
	Store *pgstore.Store
}

func New(store *pgstore.Store) *AppUC {
	return &AppUC{Store: store}
}

func (uc *AppUC) GetDocByID(ctx context.Context, id int) (*doc.Documentation, error) {
	d, err := uc.Store.Doc().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (uc *AppUC) GetAllDoc(ctx context.Context) ([]*doc.Documentation, error) {
	docs, err := uc.Store.Doc().GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

func (uc *AppUC) GetCrossed(ctx context.Context) (*crossed.Crossed, error) {
	docs, err := uc.GetAllDoc(ctx)
	if err != nil {
		return nil, err
	}

	articleNames, err := uc.Store.Article().GetAllNames(ctx)
	if err != nil {
		return nil, err
	}

	cr := make(map[string]map[string]int)
	for _, d := range docs {
		if cr[d.Name] == nil {
			cr[d.Name] = make(map[string]int)
		}

		for _, aName := range articleNames {
			cr[d.Name][aName] = 0
		}

		for _, a := range d.Articles {
			cr[d.Name][a.Name]++
		}
	}

	crsd := crossed.Crossed{
		Map:          cr,
		ArticleNames: articleNames,
	}

	return &crsd, nil
}

func (uc *AppUC) GetArticlesWithoutDoc(ctx context.Context) ([]article.Article, error) {
	arts, err := uc.Store.Article().GetWithoutDoc(ctx)
	if err != nil {
		return nil, err
	}

	return arts, nil
}
