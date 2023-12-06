package htmlview

import (
	"documentation-mini-app/internal/domain/crossed"
	"html/template"
	"io"
)

type CrossedView struct {
	t *template.Template
}

func NewCrossedView(templatePath string) (*CrossedView, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &CrossedView{t: tmpl}, nil
}

func (v *CrossedView) ToWriter(crossed crossed.Crossed, w io.Writer) error {
	err := v.t.Execute(w, crossed)
	if err != nil {
		return err
	}

	return nil
}
