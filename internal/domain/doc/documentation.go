package doc

import (
	"documentation-mini-app/internal/domain/article"
)

type Documentation struct {
	ID                       int
	Name                     string
	DefaultHighlightLanguage string
	Articles                 []article.Article
}
