package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func (v *Validator) loadSchema(serviceID string) error {
	schemaPath := v.getSchemaPath(serviceID)

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
			return nil
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

func (v *Validator) ValidateByServiceIDMap(serviceID string, config map[string]any) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return v.ValidateByServiceID(serviceID, data)
}

func (v *Validator) ReloadSchema(serviceID string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	delete(v.schemas, serviceID)

	return v.loadSchema(serviceID)
}
