package model

import (
	"strings"

	"github.com/fatih/structs"
)

// Schema is a JSON schema for a stream.
type Schema struct {
	// Type is the type of the schema, e.g. "object" - for some reason singer docs
	// have this as an array and often nullable eg: `"type": ["null", "object"]`
	Type []string `json:"type"`
	// HasAdditionalProperties indicates whether the schema allows additional properties
	// not defined in the schema.
	HasAdditionalProperties bool `json:"additionalProperties"`
	// Properties is a map of property names to their schema.
	Properties map[string]Property `json:"properties"`
}

// Property is a property in a JSON schema.
type Property struct {
	// Types is a list of types that this property can be, e.g. "string" or "integer".
	Types []string `json:"type"`
	// CustomFormat	is a custom format for this property, e.g. "date-time".
	CustomFormat string `json:"format,omitempty"`
	// For nested structures a property can have its own properties.
	Properties map[string]Property `json:"properties,omitempty"`
	// For array structures we define the type of the items in the array
	Items *ArrayItem `json:"items,omitempty"`
}

// ArrayItem is the type and properties of an item in an array.
type ArrayItem struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
}

func (s Property) IsBoolean() bool {
	return s.hasType("boolean")
}

func (s Property) IsNumber() bool {
	return s.hasType("number")
}

func (s Property) IsInteger() bool {
	return s.hasType("integer")
}

func (s Property) hasType(typeName string) bool {
	for _, t := range s.Types {
		if strings.EqualFold(t, typeName) {
			return true
		}
	}
	return false
}

func (s Property) IsDateTime() bool {
	return s.CustomFormat == "date-time"
}

// As a shortcut for simple leaf nodes we can just dump everything (we can also dump everything higher level probably too)
// If we're just going to dump everything why bother with serialisers? Good question.
// a) Initial thoughts were that we want some control on the fields we output potentially - for example
// ignoring deprecated fields.
// b) We also might need to be cleverer when it comes to catalog config that enables / disables optional fields
// in the output.
//
// Keeping this as a single callsite so it's easy to find where we're doing this in future.
func DumpToMap(input interface{}) map[string]any {
	structs.DefaultTagName = "json"
	return structs.Map(input)
}

func Optional(p Property) Property {
	for _, val := range p.Types {
		if val == "null" {
			return p
		}
	}
	p.Types = append(p.Types, "null")
	return p
}

func ArrayOf(p Property) Property {
	return Property{
		Types: []string{"array"},
		Items: &ArrayItem{
			Type:       "object",
			Properties: p.Properties,
		},
	}
}
