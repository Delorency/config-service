package schemaService

import "os"

type SchemaService struct {
	SchemaDir string
}

type SchemaServiceI interface {
	CreateSchemaFromFile(string, []byte, string) error
	GetSchema(string) ([]byte, error)
	ValidateSchema([]byte) error
	DeleteSchema(string) error
}

func NewSchemaService(schemaDir string) SchemaServiceI {
	os.MkdirAll(schemaDir, 0755)
	return &SchemaService{
		SchemaDir: schemaDir,
	}
}
