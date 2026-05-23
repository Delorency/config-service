package schemaService

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func (s *SchemaService) DeleteSchema(ctx context.Context, serviceID string) error {

	if exist := s.validator.SchemaExists(serviceID); exist {
		schemaPath := filepath.Join(s.SchemaDir, serviceID+".json")
		if err := os.Remove(schemaPath); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to delete schema file: %w", err)
			}
		}
	}

	if err := s.repo.DeleteAllServiceData(ctx, serviceID); err != nil {
		return fmt.Errorf("failed to delete service data: %w", err)
	}

	return nil
}
