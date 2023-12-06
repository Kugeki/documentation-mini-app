package pgstore

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestStore(ctx context.Context, t *testing.T, dbURL string) (
	s *Store, truncate func(ctx context.Context, tables ...string),
) {
	t.Helper()

	s, err := New(ctx, dbURL)
	if err != nil {
		t.Fatalf("open store: %v\n", err)
	}

	return s, func(ctx context.Context, tables ...string) {
		_, err := s.db.Exec(ctx, fmt.Sprintf("TRUNCATE %s RESTART IDENTITY CASCADE", strings.Join(tables, ", ")))
		defer s.Close()

		if err != nil {
			t.Errorf("truncate all: %v\n", err)
		}
	}
}
