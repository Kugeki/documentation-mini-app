package htmlview

import (
	"html/template"
	"io"
)

type TemplateView struct {
	t *template.Template
}

func New(templatePath string) (*TemplateView, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}

	return &TemplateView{t: tmpl}, nil
}

func (v *TemplateView) ToWriter(w io.Writer, data interface{}) error {
	err := v.t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
