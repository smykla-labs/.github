// Package schema provides JSON Schema generation for sync configuration.
package schema

import (
	"encoding/json"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/invopop/jsonschema"

	"github.com/smykla-labs/.github/internal/configtypes"
	"github.com/smykla-labs/.github/pkg/github"
)

// SchemaOutput represents a generated schema with its metadata.
type SchemaOutput struct {
	// Name is the short identifier for this schema (e.g., "sync-config", "settings")
	Name string
	// Filename is the output filename (e.g., "sync-config.schema.json")
	Filename string
	// Content is the generated JSON schema bytes
	Content []byte
}

// SchemaType identifies the type of schema to generate.
type SchemaType string

const (
	// SchemaSyncConfig generates schema for .github/sync-config.yml
	SchemaSyncConfig SchemaType = "sync-config"
	// SchemaSettings generates schema for .github/settings.yml
	SchemaSettings SchemaType = "settings"
)

// GenerateSchemaForType generates JSON Schema for the specified schema type.
func GenerateSchemaForType(
	modulePath, configPkgPath string,
	schemaType SchemaType,
) (*SchemaOutput, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
	}

	// Load Go comments as descriptions
	if err := reflector.AddGoComments(modulePath, configPkgPath); err != nil {
		return nil, errors.Wrap(err, "loading Go comments for schema descriptions")
	}

	var schema *jsonschema.Schema

	var output SchemaOutput

	switch schemaType {
	case SchemaSyncConfig:
		schema = reflector.Reflect(&configtypes.SyncConfig{})
		schema.ID = "https://raw.githubusercontent.com/smykla-labs/.github/main/schemas/sync-config.schema.json"
		schema.Title = "Sync Configuration"
		schema.Description = "Configuration for organization-wide label, file, and smyklot version synchronization. Place at .github/sync-config.yml in your repository."

		output.Name = "sync-config"
		output.Filename = "sync-config.schema.json"

	case SchemaSettings:
		schema = reflector.Reflect(&github.SettingsFile{})
		schema.ID = "https://raw.githubusercontent.com/smykla-labs/.github/main/schemas/settings.schema.json"
		schema.Title = "Repository Settings"
		schema.Description = "Repository settings definition for organization-wide synchronization. Place at .github/settings.yml in your repository."

		output.Name = "settings"
		output.Filename = "settings.schema.json"

	default:
		return nil, errors.Newf("unknown schema type: %s", schemaType)
	}

	schema.Version = "https://json-schema.org/draft/2020-12/schema"

	content, err := finalizeSchema(schema)
	if err != nil {
		return nil, err
	}

	output.Content = content

	return &output, nil
}

// GenerateAllSchemas generates all available schemas.
func GenerateAllSchemas(modulePath, configPkgPath string) ([]*SchemaOutput, error) {
	schemaTypes := []SchemaType{SchemaSyncConfig, SchemaSettings}
	outputs := make([]*SchemaOutput, 0, len(schemaTypes))

	for _, schemaType := range schemaTypes {
		output, err := GenerateSchemaForType(modulePath, configPkgPath, schemaType)
		if err != nil {
			return nil, errors.Wrapf(err, "generating %s schema", schemaType)
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

// finalizeSchema converts a schema to JSON and applies post-processing.
func finalizeSchema(schema *jsonschema.Schema) ([]byte, error) {
	// Convert to JSON and back to map for post-processing
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling schema to bytes")
	}

	var schemaMap map[string]any
	if err = json.Unmarshal(schemaBytes, &schemaMap); err != nil {
		return nil, errors.Wrap(err, "unmarshaling schema to map")
	}

	// Normalize descriptions (replace newlines with spaces)
	normalizeDescriptions(schemaMap)

	output, err := json.MarshalIndent(schemaMap, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "marshaling final schema")
	}

	// Add trailing newline for better git diffs
	output = append(output, '\n')

	return output, nil
}

// normalizeDescriptions recursively replaces newlines in description fields with spaces.
func normalizeDescriptions(v any) {
	switch val := v.(type) {
	case map[string]any:
		for key, value := range val {
			if key == "description" {
				if desc, ok := value.(string); ok {
					val[key] = strings.ReplaceAll(desc, "\n", " ")
				}
			} else {
				normalizeDescriptions(value)
			}
		}
	case []any:
		for _, item := range val {
			normalizeDescriptions(item)
		}
	}
}
