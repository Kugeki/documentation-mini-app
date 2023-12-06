package htmlview

import (
	"documentation-mini-app/internal/domain/article"
	"html/template"
	"io"
)

type ArticleView struct {
	t *template.Template
}

func NewArticleView(templatePath string) (*ArticleView, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &ArticleView{t: tmpl}, nil
}

func (v *ArticleView) ToWriter(art *article.Article, w io.Writer) error {
	err := v.t.Execute(w, art)
	if err != nil {
		return err
	}

	return nil
}
