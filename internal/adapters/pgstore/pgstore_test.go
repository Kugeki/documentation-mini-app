package pgstore

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func CheckTablesExistence(tables ...string) error {
	ctx := context.TODO()

	s, err := New(ctx, dbURL)
	if err != nil {
		log.Panicf("open store: %v\n", err)
	}
	defer s.Close()

	tablesString := fmt.Sprintf("'%s'", strings.Join(tables, "', '"))

	q := fmt.Sprintf("select count(*) from information_schema.tables where table_name IN (%v)", tablesString)

	var count int
	err = s.db.QueryRow(ctx, q).Scan(&count)
	if err != nil {
		return fmt.Errorf("table existence check: %v", err)
	}

	if count != len(tables) {
		return fmt.Errorf("find tables=%v, wanted=%v", count, len(tables))
	}

	return nil
}

var dbURL string

func TestMain(m *testing.M) {
	dbURL = os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		log.Panicln("need TEST_DATABASE_URL env variable")
	}

	err := CheckTablesExistence("documentation", "article", "example",
		"documentation_articles", "article_examples")
	if err != nil {
		log.Panicln(err)
	}

	os.Exit(m.Run())
}
