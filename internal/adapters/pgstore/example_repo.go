package pgstore

import (
	"context"
	"documentation-mini-app/internal/domain/example"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type ExampleRepoPG struct {
	db *pgxpool.Pool
}

func NewExampleRepoPG(db *pgxpool.Pool) *ExampleRepoPG {
	return &ExampleRepoPG{db: db}
}

func (r *ExampleRepoPG) GetByArticleID(ctx context.Context, artID int) ([]example.Example, error) {
	q := `SELECT e.id, e.name, e.description, e.code, e.output, ae.priority FROM article_examples ae
			JOIN example e on e.id = ae.example_id WHERE ae.article_id = $1`

	rows, err := r.db.Query(ctx, q, artID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]example.Example, 0)
	for rows.Next() {
		ex := example.Example{}
		err = rows.Scan(&ex.ID, &ex.Name, &ex.Description, &ex.Code, &ex.Output, &ex.Priority)
		if err != nil {
			return nil, err
		}
		res = append(res, ex)
	}

	return res, nil
}

func (r *ExampleRepoPG) GetByID(ctx context.Context, id int) (*example.Example, error) {
	q := `select e.id, e.name, e.description, e.code, e.output FROM example e where e.id = $1`

	var exa example.Example
	err := r.db.QueryRow(ctx, q, id).Scan(&exa.ID, &exa.Name, &exa.Description, &exa.Code, &exa.Output)
	if err != nil {
		return nil, err
	}

	return &exa, nil
}

func (r *ExampleRepoPG) Create(ctx context.Context, exa *example.Example) error {
	q := "insert into example(name, description, code, output) values($1, $2, $3, $4) returning id"

	var exaID int
	err := r.db.QueryRow(ctx, q, exa.Name, exa.Description, exa.Code, exa.Output).Scan(&exaID)
	exa.ID = exaID

	return err
}

func (r *ExampleRepoPG) AddToArticle(ctx context.Context, exaID int, artID int) error {
	q := "insert into article_examples(article_id, example_id) values($1, $2)"
	_, err := r.db.Exec(ctx, q, artID, exaID)
	return err
}

func (r *ExampleRepoPG) Update(ctx context.Context, exa *example.Example) error {
	q := "update example e set name = $1, description = $2, code = $3, output = $4 where e.id = $5"

	commandTag, err := r.db.Exec(ctx, q, exa.Name, exa.Description, exa.Code, exa.Output, exa.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("update example rows affected equals %v\n", commandTag.RowsAffected())
		return errors.New("update example rows affected not equals 1")
	}

	return nil
}

func (r *ExampleRepoPG) Delete(ctx context.Context, id int) error {
	q := "delete from article_examples where example_id=$1"
	_, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	q = "delete from example where id=$1"
	commandTag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("delete doc rows affected equals %v\n", commandTag.RowsAffected())
		return errors.New("example already deleted")
	}

	return nil
}
