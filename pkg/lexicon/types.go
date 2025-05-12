package lexicon

import (
	"encoding/json"
	"fmt"
	"time"
)

// Document represents a generic lexicon document
type Document struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Value     map[string]interface{} `json:"value"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

// Schema defines the structure of a lexicon type
type Schema struct {
	LexiconID string                 `json:"lexicon"`
	Revision  int                    `json:"revision"`
	Defs      map[string]interface{} `json:"defs"`
}

// Validator validates a document against its schema
type Validator interface {
	Validate(doc *Document) error
}

// SchemaValidator implements the Validator interface
type SchemaValidator struct {
	schemas map[string]*Schema
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		schemas: make(map[string]*Schema),
	}
}

// RegisterSchema registers a schema with the validator
func (v *SchemaValidator) RegisterSchema(schema *Schema) {
	v.schemas[schema.LexiconID] = schema
}

// Validate validates a document against its schema
func (v *SchemaValidator) Validate(doc *Document) error {
	schema, ok := v.schemas[doc.Type]
	if !ok {
		return fmt.Errorf("unknown document type: %s", doc.Type)
	}

	// In a real implementation, this would validate the document structure
	// against the schema definition
	return nil
}

// MarshalDocument marshals a document to JSON
func MarshalDocument(doc *Document) ([]byte, error) {
	return json.Marshal(doc)
}

// UnmarshalDocument unmarshals a document from JSON
func UnmarshalDocument(data []byte) (*Document, error) {
	var doc Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}
