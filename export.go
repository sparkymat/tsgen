package tsgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	tpl "github.com/sparkymat/tsgen/template"
)

func (s *Service) ExportToFS(path string) error {
	fileContentsMap, err := s.Export()
	if err != nil {
		return err
	}

	for filePath, fileBytes := range fileContentsMap {
		finalPath := filepath.Join(path, filePath)
		finalFolder := filepath.Dir(finalPath)

		if err := os.MkdirAll(finalFolder, 0o755); err != nil { //nolint:mnd
			return fmt.Errorf("failed to create folder '%s': %w", finalFolder, err)
		}

		if err := os.WriteFile(finalPath, fileBytes, 0o600); err != nil { //nolint:mnd
			return fmt.Errorf("failed to write file '%s': %w", finalPath, err)
		}
	}

	return nil
}

func (s *Service) Export() (map[string][]byte, error) {
	fileMap := map[string][]byte{}

	// Models
	for _, m := range s.models {
		filePath := filepath.Join("models", m.name+".ts")

		content, err := renderTemplateToString(tpl.ModelTS, m)
		if err != nil {
			return nil, fmt.Errorf("failed to export model %s: %w", m.Name, err)
		}

		fileMap[filePath] = []byte(content)
	}

	for _, thisSlice := range s.slices {
		filePath := filepath.Join("slices", thisSlice.Name+".ts")

		content, err := renderTemplateToString(tpl.SliceTS, thisSlice)
		if err != nil {
			return nil, fmt.Errorf("failed to export slice %s: %w", thisSlice.Name, err)
		}

		fileMap[filePath] = []byte(content)
	}

	return fileMap, nil
}

func renderTemplateToString(tpl string, obj any) (string, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("foo").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	if err = tmpl.Execute(&buf, obj); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}
