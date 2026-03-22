package schemaService

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *SchemaService) DeleteSchema(serviceID string) error {
	filename := filepath.Join(s.SchemaDir, serviceID+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("schema not found for service %s", serviceID)
		}
		return fmt.Errorf("failed to delete schema: %w", err)
	}
	return nil
}
