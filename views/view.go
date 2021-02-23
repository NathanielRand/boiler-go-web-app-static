package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Global variables to help us constuct our glob pattern.
const (
	LayoutDir   string = "views/layouts/"
	TemplateExt string = ".html"
)

// View struct contains one attribute "Template",
// which is a pointer to template.Template which points
// to our complied template.
type View struct {
	Template *template.Template
	Layout   string
}

// Render method on View type that renders the views.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// use a glob function to help return a slice of templates to
// include in our view.
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return files
}

// NewView makes it easier to create views.
func NewView(layout string, files ...string) *View {
	// Append passed in views with layout templates (navbar, footer, etc..)
	files = append(files, layoutFiles()...)

	// Parse appended files
	t, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Return pointer to View object with parsed files
	// as the value for Template field.
	return &View{
		Template: t,
		Layout:   layout,
	}
}
