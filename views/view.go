package views

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(templatePattern string) *TemplateRenderer {
	// test for template
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob(templatePattern)),
	}
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
