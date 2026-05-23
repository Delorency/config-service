package schemaService

import (
	"context"
	"main/internal/core/validator"
	r "main/internal/repo"
	"os"
)

type SchemaService struct {
	SchemaDir string
	repo      *r.ConfigRepository
	validator *validator.Validator
}

type SchemaServiceI interface {
	CreateSchemaFromFile(string, []byte, string) error
	GetSchema(string) ([]byte, error)
	ValidateSchema([]byte) error
	DeleteSchema(ctx context.Context, serviceID string) error
}

func NewSchemaService(schemaDir string, repo *r.ConfigRepository, validator *validator.Validator) SchemaServiceI {
	os.MkdirAll(schemaDir, 0755)
	return &SchemaService{
		SchemaDir: schemaDir,
		repo:      repo,
		validator: validator,
	}
}
