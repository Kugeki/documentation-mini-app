package pgstore

import (
	"context"
	"documentation-mini-app/internal/domain/doc"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DocRepoPG struct {
	db *pgxpool.Pool
	s  *Store
}

func NewDocRepoPG(db *pgxpool.Pool, s *Store) *DocRepoPG {
	return &DocRepoPG{db: db, s: s}
}

func (r *DocRepoPG) Create(ctx context.Context, d *doc.Documentation) error {
	q := "insert into documentation(name, default_highlight_language) values($1, $2) returning id"

	var docID int
	err := r.db.QueryRow(ctx, q, d.Name, d.DefaultHighlightLanguage).Scan(&docID)
	d.ID = docID

	if err != nil {
		return err
	}

	return nil
}

func (r *DocRepoPG) GetByID(ctx context.Context, docID int) (*doc.Documentation, error) {
	q := "select d.id, d.name, d.default_highlight_language from documentation as d where d.id = $1"

	var d doc.Documentation
	err := r.db.QueryRow(ctx, q, docID).Scan(&d.ID, &d.Name, &d.DefaultHighlightLanguage)
	if err != nil {
		return nil, err
	}

	d.Articles, err = r.s.Article().GetByDocID(ctx, docID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *DocRepoPG) GetAll(ctx context.Context) ([]*doc.Documentation, error) {
	q := "select d.id, d.name, d.default_highlight_language from documentation d"

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*doc.Documentation, 0)
	for rows.Next() {
		d := doc.Documentation{}
		err = rows.Scan(&d.ID, &d.Name, &d.DefaultHighlightLanguage)
		if err != nil {
			return nil, err
		}

		d.Articles, err = r.s.Article().GetByDocID(ctx, d.ID)
		if err != nil {
			return nil, err
		}

		res = append(res, &d)
	}

	return res, nil
}

func (r *DocRepoPG) Update(ctx context.Context, d *doc.Documentation) error {
	q := "update documentation d set name = $1, default_highlight_language = $2 where d.id = $3"

	commandTag, err := r.db.Exec(ctx, q, d.Name, d.DefaultHighlightLanguage, d.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("update doc rows affected equals %v\n", commandTag.RowsAffected())
		return errors.New("update doc rows affected not equals 1")
	}

	return nil
}

func (r *DocRepoPG) Delete(ctx context.Context, docID int) error {
	q := "delete from documentation_articles where documentation_id=$1"
	_, err := r.db.Exec(ctx, q, docID)
	if err != nil {
		return err
	}

	q = "delete from documentation where id=$1"
	commandTag, err := r.db.Exec(ctx, q, docID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("delete doc rows affected equals %v\n", commandTag.RowsAffected())
		return errors.New("delete doc rows affected not equals 1")
	}

	return nil
}
