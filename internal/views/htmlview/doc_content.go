package htmlview

import (
	"documentation-mini-app/internal/domain/doc"
	"html/template"
	"io"
)

type DocContentView struct {
	t *template.Template
}

func NewDocContentView(templatePath string) (*DocContentView, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &DocContentView{t: tmpl}, nil
}

func (v *DocContentView) ToWriter(doc []*doc.Documentation, w io.Writer) error {
	err := v.t.Execute(w, doc)
	if err != nil {
		return err
	}

	return nil
}
