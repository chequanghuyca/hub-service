package templates

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"sync"
)

//go:embed html/*.html
var templateFS embed.FS

var (
	templates     map[string]*template.Template
	templatesOnce sync.Once
)

// loadTemplates loads all HTML templates from the embedded filesystem
func loadTemplates() {
	templatesOnce.Do(func() {
		templates = make(map[string]*template.Template)

		// Load welcome template
		welcomeTmpl, err := template.ParseFS(templateFS, "html/welcome.html")
		if err != nil {
			log.Printf("Warning: Failed to load welcome.html template: %v", err)
		} else {
			templates["welcome"] = welcomeTmpl
		}

		// Load comeback template
		comebackTmpl, err := template.ParseFS(templateFS, "html/comeback.html")
		if err != nil {
			log.Printf("Warning: Failed to load comeback.html template: %v", err)
		} else {
			templates["comeback"] = comebackTmpl
		}

		log.Printf("Email templates loaded: %d templates", len(templates))
	})
}


// GetTemplate returns a parsed template by name
func GetTemplate(name string) *template.Template {
	loadTemplates()
	return templates[name]
}

// RenderTemplate renders a template with the given data and returns the HTML string
func RenderTemplate(name string, data interface{}) (string, error) {
	tmpl := GetTemplate(name)
	if tmpl == nil {
		return "", nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// HasTemplate checks if a template exists
func HasTemplate(name string) bool {
	return GetTemplate(name) != nil
}
