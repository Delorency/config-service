package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetSchemaFile(serviceID, SchemaDir string) ([]byte, error) {
	filename := filepath.Join(SchemaDir, serviceID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema not found for service %s", serviceID)
		}
		return nil, fmt.Errorf("failed to read schema: %w", err)
	}
	return data, nil
}
