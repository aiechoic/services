package openapi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenapi_NewSchema(t *testing.T) {

	type TestType struct{}

	type testCase struct {
		baseType any
		need     *Schema
	}

	testCases := []testCase{
		{
			baseType: int8(0),
			need: &Schema{
				Type:   "integer",
				Format: "int32",
			},
		},
		{
			baseType: int16(0),
			need: &Schema{
				Type:   "integer",
				Format: "int32",
			},
		},
		{
			baseType: int32(0),
			need: &Schema{
				Type:   "integer",
				Format: "int32",
			},
		},
		{
			baseType: int64(0),
			need: &Schema{
				Type:   "integer",
				Format: "int64",
			},
		},
		{
			baseType: uint8(0),
			need: &Schema{
				Type:   "integer",
				Format: "int32",
			},
		},
		{
			baseType: uint16(0),
			need: &Schema{
				Type:   "integer",
				Format: "int32",
			},
		},
		{
			baseType: uint32(0),
			need: &Schema{
				Type:   "integer",
				Format: "int32",
			},
		},
		{
			baseType: uint64(0),
			need: &Schema{
				Type:   "integer",
				Format: "int64",
			},
		},
		{
			baseType: float32(0),
			need: &Schema{
				Type:   "number",
				Format: "float",
			},
		},
		{
			baseType: float64(0),
			need: &Schema{
				Type:   "number",
				Format: "double",
			},
		},
		{
			baseType: "",
			need: &Schema{
				Type: "string",
			},
		},
		{
			baseType: true,
			need: &Schema{
				Type: "boolean",
			},
		},
		{
			baseType: []int{},
			need: &Schema{
				Type: "array",
				Items: &Schema{
					Type:   "integer",
					Format: "int32",
				},
			},
		},
		{
			baseType: TestType{},
			need: &Schema{
				Type:       "object",
				Properties: map[string]*Schema{},
			},
		},
	}

	for _, tc := range testCases {
		o := &Openapi{}
		got := o.NewSchema("", tc.baseType, ContentTypeJson)
		if got.Ref != "" {
			got = o.GetRefSchema(got.Ref)
		}
		assert.Equal(t, tc.need, got, "NewSchema(%v) = %v, want %v", tc.baseType, got, tc.need)
	}
}

func TestOpenapi_NewSchema2(t *testing.T) {
	type NestedStruct struct {
		ID   int    `json:"id" binding:"required" description:"The ID of the nested struct"`
		Name string `json:"name" description:"The name of the nested struct"`
	}
	type TestType struct {
		Name  string `json:"name" binding:"required" description:"The name of the item"`
		Email string `json:"email" binding:"required,email" description:"The email of the item"`
		NestedStruct
		Nested NestedStruct `json:"nested"`
	}

	o := &Openapi{}
	need := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {
				Type:        "string",
				Description: "The name of the item",
			},
			"email": {
				Type:        "string",
				Format:      "email",
				Description: "The email of the item",
			},
			"id": {
				Type:        "integer",
				Format:      "int32",
				Description: "The ID of the nested struct",
			},
			"nested": {
				Type: "object",
				Properties: map[string]*Schema{
					"id": {
						Type:        "integer",
						Format:      "int32",
						Description: "The ID of the nested struct",
					},
					"name": {
						Type:        "string",
						Description: "The name of the nested struct",
					},
				},
				Required: []string{"id"},
			},
		},
		Required: []string{"name", "email", "id"},
	}
	schema := o.NewSchema("", TestType{}, ContentTypeJson)
	schema = o.GetRefSchema(schema.Ref)

	assert.Equal(t, getJsonData(need), getJsonData(schema))
}

func getJsonData(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
