package jsonvalidator

import (
	"encoding/json"
	"fmt"
)

var rootSchemaPool = map[string]*RootJsonSchema{}

type RootJsonSchema struct {
	JsonSchema
	subSchemaMap map[string]*JsonSchema
}

func NewRootJsonSchema(bytes []byte) (*RootJsonSchema, error) {
	var rootSchema *RootJsonSchema

	// Check if the string s is a valid json.
	err := json.Unmarshal(bytes, &rootSchema)
	if err != nil {
		return nil, err
	}

	rootSchema.subSchemaMap = make(map[string]*JsonSchema)

	if rootSchema.Id != nil {
		if _, ok := rootSchemaPool[string(*rootSchema.Id)]; !ok {
			rootSchemaPool[string(*rootSchema.Id)] = rootSchema
		}
	}

	err = rootSchema.scanSchema("", string(*rootSchema.Id))
	if err != nil {
		fmt.Println("[JsonSchema DEBUG] scanSchema() " +
			"failed: " + err.Error())
		return nil, err
	}

	return rootSchema, nil
}

func (rs *RootJsonSchema) validate(bytes []byte) (bool, error) {
	var id string
	if rs.Id != nil {
		id = string(*rs.Id)
	} else {
		id = ""
	}

	return rs.validateJsonData("", bytes, id)
}
