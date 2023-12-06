package exampleuc

import (
	"context"
	"documentation-mini-app/internal/adapters/pgstore"
	"documentation-mini-app/internal/domain/example"
)

type ExampleUC struct {
	Store *pgstore.Store
}

func New(store *pgstore.Store) *ExampleUC {
	return &ExampleUC{Store: store}
}

func (uc *ExampleUC) GetExampleByID(ctx context.Context, id int) (*example.Example, error) {
	exa, err := uc.Store.Example().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return exa, nil
}

func (uc *ExampleUC) CreateExample(ctx context.Context, exa *example.Example) error {
	err := uc.Store.Example().Create(ctx, exa)
	return err
}

func (uc *ExampleUC) AddExampleToArticle(ctx context.Context, exaID int, artID int) error {
	if artID != 0 {
		err := uc.Store.Example().AddToArticle(ctx, exaID, artID)
		return err
	}
	return nil
}

func (uc *ExampleUC) UpdateExample(ctx context.Context, exa *example.Example) error {
	err := uc.Store.Example().Update(ctx, exa)
	return err
}

func (uc *ExampleUC) DeleteExample(ctx context.Context, id int) error {
	err := uc.Store.Example().Delete(ctx, id)
	return err
}
