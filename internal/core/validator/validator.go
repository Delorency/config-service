package validator

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type Validator struct {
	schemaLoader gojsonschema.JSONLoader
}

func NewValidator(schemaPath string) (*Validator, error) {
	loader := gojsonschema.NewReferenceLoader("file://" + schemaPath)

	_, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	return &Validator{
		schemaLoader:  loader,
		usePerService: false,
	}, nil
}
