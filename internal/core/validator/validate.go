package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaExists проверяет существование файла схемы
func (v *Validator) SchemaExists(serviceID string) bool {
	schemaPath := v.getSchemaPath(serviceID)
	_, err := os.Stat(schemaPath)
	return !os.IsNotExist(err)
}

func (v *Validator) loadSchema(serviceID string) error {
	schemaPath := v.getSchemaPath(serviceID)

	// Проверяем существование файла
	if !v.SchemaExists(serviceID) {
		return fmt.Errorf("schema file not found for service %s at path: %s", serviceID, schemaPath)
	}

	loader := gojsonschema.NewReferenceLoader("file://" + schemaPath)

	if _, err := gojsonschema.NewSchema(loader); err != nil {
		return fmt.Errorf("invalid schema for service %s: %w", serviceID, err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	v.schemas[serviceID] = loader

	return nil
}

func (v *Validator) ValidateByServiceID(serviceID string, config []byte) error {
	v.mu.RLock()
	loader, ok := v.schemas[serviceID]
	v.mu.RUnlock()

	if !ok {
		if err := v.loadSchema(serviceID); err != nil {
			return fmt.Errorf("failed to load schema for service %s: %w", serviceID, err)
		}

		v.mu.RLock()
		loader = v.schemas[serviceID]
		v.mu.RUnlock()
	}

	documentLoader := gojsonschema.NewBytesLoader(config)

	result, err := gojsonschema.Validate(loader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var sb strings.Builder
		for _, err := range result.Errors() {
			sb.WriteString(fmt.Sprintf("- %s\n", err))
		}
		return fmt.Errorf("schema validation failed for service %s:\n%s", serviceID, sb.String())
	}

	return nil
}

// ValidateByServiceIDWithExistsCheck проверяет существование схемы И валидирует конфиг
func (v *Validator) ValidateByServiceIDWithExistsCheck(serviceID string, config []byte) error {
	// Сначала проверяем существование файла схемы
	if !v.SchemaExists(serviceID) {
		return fmt.Errorf("schema not found for service '%s'. Please upload schema first", serviceID)
	}

	// Затем валидируем конфиг
	return v.ValidateByServiceID(serviceID, config)
}

func (v *Validator) ValidateByServiceIDMap(serviceID string, config map[string]any) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return v.ValidateByServiceID(serviceID, data)
}

func (v *Validator) ReloadSchema(serviceID string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	delete(v.schemas, serviceID)

	return v.loadSchema(serviceID)
}
