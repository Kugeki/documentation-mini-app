package article

import (
	"documentation-mini-app/internal/domain/example"
)

type Article struct {
	ID          int
	Name        string
	Description string
	Examples    []example.Example
}
