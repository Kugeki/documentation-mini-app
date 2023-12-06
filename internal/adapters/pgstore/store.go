package pgstore

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Store struct {
	db          *pgxpool.Pool
	docRepo     *DocRepoPG
	articleRepo *ArticleRepoPG
	exampleRepo *ExampleRepoPG
}

// New connects database. Need call Close after this.
func New(ctx context.Context, dbURL string) (*Store, error) {
	s := &Store{}

	err := s.open(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// open connects database. Need call Close after this.
func (s *Store) open(ctx context.Context, dbURL string) error {
	if s.db != nil {
		log.Println("trying to open store that not closed")
		s.Close()
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("pool create: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return fmt.Errorf("db ping: %w", err)
	}

	s.db = pool

	return nil
}

func (s *Store) Close() {
	if s.db == nil {
		log.Println("trying to close nil db")
		return
	}

	s.db.Close()
}

func (s *Store) Doc() *DocRepoPG {
	if s.docRepo == nil {
		s.docRepo = NewDocRepoPG(s.db, s)
	}

	return s.docRepo
}

func (s *Store) Article() *ArticleRepoPG {
	if s.articleRepo == nil {
		s.articleRepo = NewArticleRepoPG(s.db, s)
	}

	return s.articleRepo
}

func (s *Store) Example() *ExampleRepoPG {
	if s.exampleRepo == nil {
		s.exampleRepo = NewExampleRepoPG(s.db)
	}

	return s.exampleRepo
}
