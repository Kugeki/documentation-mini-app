package docuc

import (
	"context"
	"documentation-mini-app/internal/adapters/pgstore"
	"documentation-mini-app/internal/domain/doc"
)

type DocUC struct {
	Store *pgstore.Store
}

func New(store *pgstore.Store) *DocUC {
	return &DocUC{Store: store}
}

func (uc *DocUC) GetDocByID(ctx context.Context, docID int) (*doc.Documentation, error) {
	d, err := uc.Store.Doc().GetByID(ctx, docID)
	return d, err
}

func (uc *DocUC) CreateDoc(ctx context.Context, d *doc.Documentation) error {
	err := uc.Store.Doc().Create(ctx, d)
	return err
}

func (uc *DocUC) UpdateDoc(ctx context.Context, d *doc.Documentation) error {
	err := uc.Store.Doc().Update(ctx, d)
	return err
}

func (uc *DocUC) DeleteDoc(ctx context.Context, docID int) error {
	err := uc.Store.Doc().Delete(ctx, docID)
	return err
}
