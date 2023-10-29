package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"text/template"
)

type Generator struct {
	filename string
	rendered []byte
}

func NewGenerator(filename string) *Generator {
	return &Generator{
		filename: filename,
	}
}

func parseTemplate(body string, data any) ([]byte, error) {
	// Parse template
	tmpl, err := template.New("gen").
		Funcs(
			template.FuncMap{
				"ToTitle": cases.Title(language.English).String,
			}).
		Parse(body)
	if err != nil {
		return nil, fmt.Errorf("failed during parsing template: %w", err)
	}

	// Execute template
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return nil, fmt.Errorf("failed during executing template: %w", err)
	}

	return buffer.Bytes(), nil
}

func (g *Generator) Render(body string, data any) error {
	// Render template
	rendered, err := parseTemplate(body, data)
	if err != nil {
		return err
	}

	// Format template with gofmt
	formatted, err := format.Source(rendered)
	if err != nil {
		return fmt.Errorf("failed during formatting:\n%s\n%w", rendered, err)
	}
	g.rendered = formatted

	return err
}

func (g *Generator) Write() error {
	f, err := os.OpenFile(g.filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed opening file %s\n  %w", g.filename, err)
	}
	defer f.Close()

	if _, err = f.Write(g.rendered); err != nil {
		return fmt.Errorf("failed writing file %s\n  %w", g.filename, err)
	}

	return nil
}
