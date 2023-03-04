package reflectopenapi_test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	reflectopenapi "github.com/podhmo/reflect-openapi"
)

// This is Owner of something
type Owner struct {
	// name of owner
	Name string `json:"name" openapi-override:"{'pattern': '^[A-Z][-A-Za-z]+$'}"`
	Age  int    `json:"age,omitempty"` // age of owner
}

// input parameters
type ListOwnerInput struct {
	// sort option
	Sort string `json:"sort" in:"query" openapi-override:"{'enum': ['desc', 'asc'], 'default': 'asc'}"`
}

// Returns list of owners.
func ListOwner(ctx context.Context, input ListOwnerInput) ([]*Owner, error) {
	return nil, nil
}

func ExampleConfig() {
	c := reflectopenapi.Config{
		TagNameOption: &reflectopenapi.TagNameOption{
			NameTag:        "json",
			ParamTypeTag:   "in",
			DescriptionTag: "description",
			OverrideTag:    "openapi-override",
		},
		SkipValidation: true,
		EnableAutoTag:  true,
		Extractor:      shapeCfg,
	}
	doc, _ := c.BuildDoc(context.Background(), func(m *reflectopenapi.Manager) {
		m.RegisterFunc(ListOwner).After(func(op *openapi3.Operation) {
			m.Doc.AddOperation("/owners", "GET", op)
		})
	})
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "@@")
	enc.Encode(doc)
	// Output:
	// {
	// @@"components": {
	// @@@@"schemas": {
	// @@@@@@"Owner": {
	// @@@@@@@@"description": "This is Owner of something",
	// @@@@@@@@"properties": {
	// @@@@@@@@@@"age": {
	// @@@@@@@@@@@@"description": "age of owner",
	// @@@@@@@@@@@@"type": "integer"
	// @@@@@@@@@@},
	// @@@@@@@@@@"name": {
	// @@@@@@@@@@@@"description": "name of owner",
	// @@@@@@@@@@@@"pattern": "^[A-Z][-A-Za-z]+$",
	// @@@@@@@@@@@@"type": "string"
	// @@@@@@@@@@}
	// @@@@@@@@},
	// @@@@@@@@"required": [
	// @@@@@@@@@@"name"
	// @@@@@@@@],
	// @@@@@@@@"title": "Owner",
	// @@@@@@@@"type": "object"
	// @@@@@@}
	// @@@@}
	// @@},
	// @@"info": {
	// @@@@"title": "Sample API",
	// @@@@"version": "0.0.0"
	// @@},
	// @@"openapi": "3.0.0",
	// @@"paths": {
	// @@@@"/owners": {
	// @@@@@@"get": {
	// @@@@@@@@"description": "Returns list of owners.",
	// @@@@@@@@"operationId": "github.com/podhmo/reflect-openapi_test.ListOwner",
	// @@@@@@@@"parameters": [
	// @@@@@@@@@@{
	// @@@@@@@@@@@@"description": "sort option",
	// @@@@@@@@@@@@"in": "query",
	// @@@@@@@@@@@@"name": "sort",
	// @@@@@@@@@@@@"schema": {
	// @@@@@@@@@@@@@@"default": "asc",
	// @@@@@@@@@@@@@@"enum": [
	// @@@@@@@@@@@@@@@@"desc",
	// @@@@@@@@@@@@@@@@"asc"
	// @@@@@@@@@@@@@@],
	// @@@@@@@@@@@@@@"type": "string"
	// @@@@@@@@@@@@}
	// @@@@@@@@@@}
	// @@@@@@@@],
	// @@@@@@@@"responses": {
	// @@@@@@@@@@"200": {
	// @@@@@@@@@@@@"content": {
	// @@@@@@@@@@@@@@"application/json": {
	// @@@@@@@@@@@@@@@@"schema": {
	// @@@@@@@@@@@@@@@@@@"items": {
	// @@@@@@@@@@@@@@@@@@@@"$ref": "#/components/schemas/Owner"
	// @@@@@@@@@@@@@@@@@@},
	// @@@@@@@@@@@@@@@@@@"type": "array"
	// @@@@@@@@@@@@@@@@}
	// @@@@@@@@@@@@@@}
	// @@@@@@@@@@@@},
	// @@@@@@@@@@@@"description": ""
	// @@@@@@@@@@},
	// @@@@@@@@@@"default": {
	// @@@@@@@@@@@@"description": ""
	// @@@@@@@@@@}
	// @@@@@@@@},
	// @@@@@@@@"summary": "Returns list of owners.",
	// @@@@@@@@"tags": [
	// @@@@@@@@@@"reflect-openapi_test"
	// @@@@@@@@]
	// @@@@@@}
	// @@@@}
	// @@},
	// @@"servers": [
	// @@@@{
	// @@@@@@"description": "local development server",
	// @@@@@@"url": "http://localhost:8888"
	// @@@@}
	// @@],
	// @@"tags": [
	// @@@@{
	// @@@@@@"name": "reflect-openapi_test"
	// @@@@}
	// @@]
	// }
}
