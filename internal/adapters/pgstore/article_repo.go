package pgstore

import (
	"context"
	"documentation-mini-app/internal/domain/article"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type ArticleRepoPG struct {
	db *pgxpool.Pool
	s  *Store
}

func NewArticleRepoPG(db *pgxpool.Pool, s *Store) *ArticleRepoPG {
	return &ArticleRepoPG{db: db, s: s}
}

func (r *ArticleRepoPG) Create(ctx context.Context, art *article.Article) error {
	q := "insert into article(name, description) values($1, $2) returning id"

	var artID int
	err := r.db.QueryRow(ctx, q, art.Name, art.Description).Scan(&artID)
	art.ID = artID

	return err
}

func (r *ArticleRepoPG) GetByID(ctx context.Context, id int) (*article.Article, error) {
	q := `SELECT a.id, a.name, a.description FROM article a WHERE id = $1`

	var art article.Article
	err := r.db.QueryRow(ctx, q, id).Scan(&art.ID, &art.Name, &art.Description)
	if err != nil {
		return nil, err
	}

	art.Examples, err = r.s.Example().GetByArticleID(ctx, art.ID)
	if err != nil {
		return nil, err
	}

	return &art, nil
}

func (r *ArticleRepoPG) GetAllNames(ctx context.Context) ([]string, error) {
	q := "SELECT a.name FROM article a ORDER BY a.name"

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		res = append(res, name)
	}

	return res, nil
}

func (r *ArticleRepoPG) GetByDocID(ctx context.Context, docID int) ([]article.Article, error) {
	q := `SELECT a.id, a.name, a.description FROM documentation_articles da 
			JOIN article a on a.id = da.article_id WHERE documentation_id = $1`

	rows, err := r.db.Query(ctx, q, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]article.Article, 0)
	for rows.Next() {
		art := article.Article{}
		err = rows.Scan(&art.ID, &art.Name, &art.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, art)
	}

	return res, nil
}

func (r *ArticleRepoPG) AddToDoc(ctx context.Context, artID int, docID int) error {
	q := "insert into documentation_articles(documentation_id, article_id) values($1, $2)"
	_, err := r.db.Exec(ctx, q, docID, artID)
	return err
}

func (r *ArticleRepoPG) Update(ctx context.Context, art *article.Article) error {
	q := "update article a set name = $1, description = $2 where a.id = $3"

	commandTag, err := r.db.Exec(ctx, q, art.Name, art.Description, art.ID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("update article rows affected equals %v\n", commandTag.RowsAffected())
		return errors.New("update article rows affected not equals 1")
	}

	return nil
}

func (r *ArticleRepoPG) Delete(ctx context.Context, artID int) error {
	q := "delete from documentation_articles where article_id=$1"
	_, err := r.db.Exec(ctx, q, artID)
	if err != nil {
		return err
	}

	q = "delete from article where id=$1"
	commandTag, err := r.db.Exec(ctx, q, artID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("delete doc rows affected equals %v\n", commandTag.RowsAffected())
		return errors.New("article already deleted")
	}

	return nil
}

func (r *ArticleRepoPG) GetWithoutDoc(ctx context.Context) ([]article.Article, error) {
	q := `
	select a.id, a.name, a.description from article a 
    left join documentation_articles da ON da.article_id = a.id
    where da.documentation_id is null
    `

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]article.Article, 0)
	for rows.Next() {
		art := article.Article{}
		err = rows.Scan(&art.ID, &art.Name, &art.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, art)
	}

	return res, nil
}
