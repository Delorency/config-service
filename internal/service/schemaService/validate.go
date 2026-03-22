package schemaService

import (
	"encoding/json"
	"fmt"
)

func (s *SchemaService) ValidateSchema(schemaData []byte) error {
	var schema map[string]any
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if _, ok := schema["type"]; !ok {
		return fmt.Errorf("schema must have 'type' field")
	}

	return nil
}
