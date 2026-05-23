package schemaService

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (s *SchemaService) CreateSchemaFromFile(serviceID string, fileContent []byte, filename string) error {
	if serviceID == "" {
		return fmt.Errorf("service_id is required")
	}

	if ext := filepath.Ext(filename); ext != ".json" {
		return fmt.Errorf("only JSON files are allowed, got %s", ext)
	}

	var schema map[string]any
	if err := json.Unmarshal(fileContent, &schema); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if _, ok := schema["type"]; !ok {
		return fmt.Errorf("schema must have 'type' field")
	}

	savePath := filepath.Join(s.SchemaDir, serviceID+".json")
	if err := os.WriteFile(savePath, fileContent, 0644); err != nil {
		return fmt.Errorf("failed to save schema: %w", err)
	}

	return nil
}
