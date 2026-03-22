package validator

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func (v *Validator) Validate(config []byte) error {
	documentLoader := gojsonschema.NewBytesLoader(config)

	result, err := gojsonschema.Validate(v.schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		errors := ""
		for _, err := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", err)
		}
		return fmt.Errorf("schema validation failed:\n%s", errors)
	}

	return nil
}

func (v *Validator) ValidateMap(config map[string]any) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return v.Validate(data)
}
