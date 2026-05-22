package schemaService

import (
	"fmt"
	"main/internal/utils"
	"os"
	"path/filepath"
)

func (s *SchemaService) GetSchema(serviceID string) ([]byte, error) {
	return utils.GetSchemaFile(serviceID, s.SchemaDir)
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
