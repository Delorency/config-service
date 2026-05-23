package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	schemaDir string
	schemas   map[string]gojsonschema.JSONLoader
	mu        sync.RWMutex
}

func NewValidator(schemaDir string) (*Validator, error) {
	if err := os.MkdirAll(schemaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create schema directory: %w", err)
	}

	return &Validator{
		schemaDir: schemaDir,
		schemas:   make(map[string]gojsonschema.JSONLoader),
	}, nil
}

func (v *Validator) getSchemaPath(serviceID string) string {
	return filepath.Join(v.schemaDir, serviceID+".json")
}
