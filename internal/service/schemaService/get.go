package schemaService

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *SchemaService) GetSchema(serviceID string) ([]byte, error) {
	filename := filepath.Join(s.SchemaDir, serviceID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema not found for service %s", serviceID)
		}
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}
	return data, nil
}

func (s *SchemaService) ListSchemas() ([]map[string]string, error) {
	files, err := filepath.Glob(filepath.Join(s.SchemaDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}

	schemas := make([]map[string]string, 0, len(files))
	for _, file := range files {
		info, _ := os.Stat(file)
		serviceID := filepath.Base(file[:len(file)-5])

		schemas = append(schemas, map[string]string{
			"service_id": serviceID,
			"file":       filepath.Base(file),
			"size":       fmt.Sprintf("%d bytes", info.Size()),
			"modified":   info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	return schemas, nil
}
