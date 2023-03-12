package goecs

import "github.com/xeipuuv/gojsonschema"

const BASE_COMPONENT_SCHEMA = `
{
	"$schema": "http://json-schema.org/draft-07/schema#",
	"type": "object",
	"properties": {
	  "id": {
		"type": "object",
		"properties": {
		  "namespace": {
			"type": "string"
		  },
		  "path": {
			"type": "string"
		  }
		},
		"required": ["namespace", "path"]
	  },
	  "data": {
		"oneOf": [
		  {
			"type": "object"
		  },
		  {
			"type": "array"
		  }
		]
	  }
	},
	"required": ["id", "data"]
  }
`

func ValidateComponentFromBaseSchema(component Component) (bool, error) {
	schemaLoader := gojsonschema.NewStringLoader(BASE_COMPONENT_SCHEMA)

	componentLoader := gojsonschema.NewGoLoader(map[string]interface{}{
		"id":   component.Id,
		"data": component.Data,
	})

	result, err := gojsonschema.Validate(schemaLoader, componentLoader)
	if err != nil {
		return false, err
	}

	return result.Valid(), nil
}

func ValidateComponentFromSchema(component Component, schema string) (bool, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema)

	componentLoader := gojsonschema.NewGoLoader(map[string]interface{}{
		"id":   component.Id,
		"data": component.Data,
	})

	result, err := gojsonschema.Validate(schemaLoader, componentLoader)
	if err != nil {
		return false, err
	}

	return result.Valid(), nil
}
