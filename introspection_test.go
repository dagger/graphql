package graphql_test

import (
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/location"
	"github.com/dagger/graphql/testutil"
)

func g(t *testing.T, p graphql.Params) *graphql.Result {
	return graphql.Do(p)
}

func TestIntrospection_ExecutesAnIntrospectionQuery(t *testing.T) {
	emptySchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "QueryRoot",
			Fields: graphql.Fields{
				"onlyField": &graphql.Field{
					Type: graphql.String,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	expectedDataSubSet := map[string]any{
		"__schema": map[string]any{
			"mutationType":     nil,
			"subscriptionType": nil,
			"queryType": map[string]any{
				"name": "QueryRoot",
			},
			"types": []any{
				map[string]any{
					"kind":          "OBJECT",
					"name":          "QueryRoot",
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind": "OBJECT",
					"name": "__Schema",
					"fields": []any{
						map[string]any{
							"name": "types",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]any{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]any{
											"kind": "OBJECT",
											"name": "__Type",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "queryType",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind": "OBJECT",
									"name": "__Type",
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "mutationType",
							"args": []any{},
							"type": map[string]any{
								"kind": "OBJECT",
								"name": "__Type",
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "subscriptionType",
							"args": []any{},
							"type": map[string]any{
								"kind": "OBJECT",
								"name": "__Type",
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "directives",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]any{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]any{
											"kind": "OBJECT",
											"name": "__Directive",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind": "OBJECT",
					"name": "__Type",
					"fields": []any{
						map[string]any{
							"name": "kind",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "ENUM",
									"name":   "__TypeKind",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "name",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "description",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "fields",
							"args": []any{
								map[string]any{
									"name": "includeDeprecated",
									"type": map[string]any{
										"kind":   "SCALAR",
										"name":   "Boolean",
										"ofType": nil,
									},
									"defaultValue": "false",
								},
							},
							"type": map[string]any{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]any{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]any{
										"kind":   "OBJECT",
										"name":   "__Field",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "interfaces",
							"args": []any{},
							"type": map[string]any{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]any{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]any{
										"kind":   "OBJECT",
										"name":   "__Type",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "possibleTypes",
							"args": []any{},
							"type": map[string]any{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]any{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]any{
										"kind":   "OBJECT",
										"name":   "__Type",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "enumValues",
							"args": []any{
								map[string]any{
									"name": "includeDeprecated",
									"type": map[string]any{
										"kind":   "SCALAR",
										"name":   "Boolean",
										"ofType": nil,
									},
									"defaultValue": "false",
								},
							},
							"type": map[string]any{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]any{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]any{
										"kind":   "OBJECT",
										"name":   "__EnumValue",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "inputFields",
							"args": []any{},
							"type": map[string]any{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]any{
									"kind": "NON_NULL",
									"name": nil,
									"ofType": map[string]any{
										"kind":   "OBJECT",
										"name":   "__InputValue",
										"ofType": nil,
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "ofType",
							"args": []any{},
							"type": map[string]any{
								"kind":   "OBJECT",
								"name":   "__Type",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind":        "ENUM",
					"name":        "__TypeKind",
					"fields":      nil,
					"inputFields": nil,
					"interfaces":  nil,
					"enumValues": []any{
						map[string]any{
							"name":              "SCALAR",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "OBJECT",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "INTERFACE",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "UNION",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "ENUM",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "INPUT_OBJECT",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "LIST",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "NON_NULL",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"possibleTypes": nil,
				},
				map[string]any{
					"kind":          "SCALAR",
					"name":          "String",
					"fields":        nil,
					"inputFields":   nil,
					"interfaces":    nil,
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind":          "SCALAR",
					"name":          "Boolean",
					"fields":        nil,
					"inputFields":   nil,
					"interfaces":    nil,
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind": "OBJECT",
					"name": "__Field",
					"fields": []any{
						map[string]any{
							"name": "name",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "description",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "args",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]any{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]any{
											"kind": "OBJECT",
											"name": "__InputValue",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "type",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "OBJECT",
									"name":   "__Type",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "isDeprecated",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "deprecationReason",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind": "OBJECT",
					"name": "__InputValue",
					"fields": []any{
						map[string]any{
							"name": "name",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "description",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "type",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "OBJECT",
									"name":   "__Type",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "defaultValue",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind": "OBJECT",
					"name": "__EnumValue",
					"fields": []any{
						map[string]any{
							"name": "name",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "description",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "isDeprecated",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "deprecationReason",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind": "OBJECT",
					"name": "__Directive",
					"fields": []any{
						map[string]any{
							"name": "name",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "description",
							"args": []any{},
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "locations",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]any{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]any{
											"kind": "ENUM",
											"name": "__DirectiveLocation",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "args",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind": "LIST",
									"name": nil,
									"ofType": map[string]any{
										"kind": "NON_NULL",
										"name": nil,
										"ofType": map[string]any{
											"kind": "OBJECT",
											"name": "__InputValue",
										},
									},
								},
							},
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name": "onOperation",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      true,
							"deprecationReason": "Use `locations`.",
						},
						map[string]any{
							"name": "onFragment",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      true,
							"deprecationReason": "Use `locations`.",
						},
						map[string]any{
							"name": "onField",
							"args": []any{},
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
							"isDeprecated":      true,
							"deprecationReason": "Use `locations`.",
						},
					},
					"inputFields":   nil,
					"interfaces":    []any{},
					"enumValues":    nil,
					"possibleTypes": nil,
				},
				map[string]any{
					"kind":        "ENUM",
					"name":        "__DirectiveLocation",
					"fields":      nil,
					"inputFields": nil,
					"interfaces":  nil,
					"enumValues": []any{
						map[string]any{
							"name":              "QUERY",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "MUTATION",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "SUBSCRIPTION",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "FIELD",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "FRAGMENT_DEFINITION",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "FRAGMENT_SPREAD",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
						map[string]any{
							"name":              "INLINE_FRAGMENT",
							"isDeprecated":      false,
							"deprecationReason": nil,
						},
					},
					"possibleTypes": nil,
				},
			},
			"directives": []any{
				map[string]any{
					"name": "include",
					"locations": []any{
						"FIELD",
						"FRAGMENT_SPREAD",
						"INLINE_FRAGMENT",
					},
					"args": []any{
						map[string]any{
							"defaultValue": nil,
							"name":         "if",
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
						},
					},
					// deprecated, but included for coverage till removed
					"onOperation": false,
					"onFragment":  true,
					"onField":     true,
				},
				map[string]any{
					"name": "skip",
					"locations": []any{
						"FIELD",
						"FRAGMENT_SPREAD",
						"INLINE_FRAGMENT",
					},
					"args": []any{
						map[string]any{
							"defaultValue": nil,
							"name":         "if",
							"type": map[string]any{
								"kind": "NON_NULL",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "Boolean",
									"ofType": nil,
								},
							},
						},
					},
					// deprecated, but included for coverage till removed
					"onOperation": false,
					"onFragment":  true,
					"onField":     true,
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        emptySchema,
		RequestString: testutil.IntrospectionQuery,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expectedDataSubSet) {
		t.Fatalf("unexpected, result does not contain subset of expected data")
	}
}

func TestIntrospection_ExecutesAnInputObject(t *testing.T) {

	testInputObject := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "TestInputObject",
		Fields: graphql.InputObjectConfigFieldMap{
			"a": &graphql.InputObjectFieldConfig{
				Type:         graphql.String,
				DefaultValue: "foo",
			},
			"b": &graphql.InputObjectFieldConfig{
				Type: graphql.NewList(graphql.String),
			},
		},
	})
	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"field": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					&graphql.ArgumentConfig{
						Name: "complex",
						Type: testInputObject,
					},
				},
				Resolve: func(p graphql.ResolveParams) (any, error) {
					return p.Args["complex"], nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __schema {
          types {
            kind
            name
            inputFields {
              name
              type { ...TypeRef }
              defaultValue
            }
          }
        }
      }

      fragment TypeRef on __Type {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
            }
          }
        }
      }
    `
	expectedDataSubSet := map[string]any{
		"__schema": map[string]any{
			"types": []any{
				map[string]any{
					"kind": "INPUT_OBJECT",
					"name": "TestInputObject",
					"inputFields": []any{
						map[string]any{
							"name": "a",
							"type": map[string]any{
								"kind":   "SCALAR",
								"name":   "String",
								"ofType": nil,
							},
							"defaultValue": `"foo"`,
						},
						map[string]any{
							"name": "b",
							"type": map[string]any{
								"kind": "LIST",
								"name": nil,
								"ofType": map[string]any{
									"kind":   "SCALAR",
									"name":   "String",
									"ofType": nil,
								},
							},
							"defaultValue": nil,
						},
					},
				},
			},
		},
	}

	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expectedDataSubSet) {
		t.Fatalf("unexpected, result does not contain subset of expected data")
	}
}

func TestIntrospection_SupportsThe__TypeRootField(t *testing.T) {

	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"testField": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestType") {
          name
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]any{
			"__type": map[string]any{
				"name": "TestType",
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestIntrospection_IdentifiesDeprecatedFields(t *testing.T) {

	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"nonDeprecated": &graphql.Field{
				Type: graphql.String,
			},
			"deprecated": &graphql.Field{
				Type:              graphql.String,
				DeprecationReason: "Removed in 1.0",
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestType") {
          name
          fields(includeDeprecated: true) {
            name
            isDeprecated,
            deprecationReason
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]any{
			"__type": map[string]any{
				"name": "TestType",
				"fields": []any{
					map[string]any{
						"name":              "nonDeprecated",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
					map[string]any{
						"name":              "deprecated",
						"isDeprecated":      true,
						"deprecationReason": "Removed in 1.0",
					},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expected.Data.(map[string]any)) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestIntrospection_RespectsTheIncludeDeprecatedParameterForFields(t *testing.T) {

	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"nonDeprecated": &graphql.Field{
				Type: graphql.String,
			},
			"deprecated": &graphql.Field{
				Type:              graphql.String,
				DeprecationReason: "Removed in 1.0",
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestType") {
          name
          trueFields: fields(includeDeprecated: true) {
            name
          }
          falseFields: fields(includeDeprecated: false) {
            name
          }
          omittedFields: fields {
            name
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]any{
			"__type": map[string]any{
				"name": "TestType",
				"trueFields": []any{
					map[string]any{
						"name": "nonDeprecated",
					},
					map[string]any{
						"name": "deprecated",
					},
				},
				"falseFields": []any{
					map[string]any{
						"name": "nonDeprecated",
					},
				},
				"omittedFields": []any{
					map[string]any{
						"name": "nonDeprecated",
					},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expected.Data.(map[string]any)) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestIntrospection_IdentifiesDeprecatedEnumValues(t *testing.T) {

	testEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"NONDEPRECATED": &graphql.EnumValueConfig{
				Value: 0,
			},
			"DEPRECATED": &graphql.EnumValueConfig{
				Value:             1,
				DeprecationReason: "Removed in 1.0",
			},
			"ALSONONDEPRECATED": &graphql.EnumValueConfig{
				Value: 2,
			},
		},
	})
	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"testEnum": &graphql.Field{
				Type: testEnum,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestEnum") {
          name
          enumValues(includeDeprecated: true) {
            name
            isDeprecated,
            deprecationReason
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]any{
			"__type": map[string]any{
				"name": "TestEnum",
				"enumValues": []any{
					map[string]any{
						"name":              "NONDEPRECATED",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
					map[string]any{
						"name":              "DEPRECATED",
						"isDeprecated":      true,
						"deprecationReason": "Removed in 1.0",
					},
					map[string]any{
						"name":              "ALSONONDEPRECATED",
						"isDeprecated":      false,
						"deprecationReason": nil,
					},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expected.Data.(map[string]any)) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestIntrospection_RespectsTheIncludeDeprecatedParameterForEnumValues(t *testing.T) {

	testEnum := graphql.NewEnum(graphql.EnumConfig{
		Name: "TestEnum",
		Values: graphql.EnumValueConfigMap{
			"NONDEPRECATED": &graphql.EnumValueConfig{
				Value: 0,
			},
			"DEPRECATED": &graphql.EnumValueConfig{
				Value:             1,
				DeprecationReason: "Removed in 1.0",
			},
			"ALSONONDEPRECATED": &graphql.EnumValueConfig{
				Value: 2,
			},
		},
	})
	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"testEnum": &graphql.Field{
				Type: testEnum,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type(name: "TestEnum") {
          name
          trueValues: enumValues(includeDeprecated: true) {
            name
          }
          falseValues: enumValues(includeDeprecated: false) {
            name
          }
          omittedValues: enumValues {
            name
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]any{
			"__type": map[string]any{
				"name": "TestEnum",
				"trueValues": []any{
					map[string]any{
						"name": "NONDEPRECATED",
					},
					map[string]any{
						"name": "DEPRECATED",
					},
					map[string]any{
						"name": "ALSONONDEPRECATED",
					},
				},
				"falseValues": []any{
					map[string]any{
						"name": "NONDEPRECATED",
					},
					map[string]any{
						"name": "ALSONONDEPRECATED",
					},
				},
				"omittedValues": []any{
					map[string]any{
						"name": "NONDEPRECATED",
					},
					map[string]any{
						"name": "ALSONONDEPRECATED",
					},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expected.Data.(map[string]any)) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestIntrospection_FailsAsExpectedOnThe__TypeRootFieldWithoutAnArg(t *testing.T) {

	testType := graphql.NewObject(graphql.ObjectConfig{
		Name: "TestType",
		Fields: graphql.Fields{
			"testField": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: testType,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        __type {
          name
        }
      }
    `
	expected := &graphql.Result{
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Field "__type" argument "name" of type "String!" ` +
					`is required but not provided.`,
				Locations: []location.SourceLocation{
					{Line: 3, Column: 9},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestIntrospection_ExposesDescriptionsOnTypesAndFields(t *testing.T) {

	queryRoot := graphql.NewObject(graphql.ObjectConfig{
		Name: "QueryRoot",
		Fields: graphql.Fields{
			"onlyField": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryRoot,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        schemaType: __type(name: "__Schema") {
          name,
          description,
          fields {
            name,
            description
          }
        }
      }
    `

	expected := &graphql.Result{
		Data: map[string]any{
			"schemaType": map[string]any{
				"name": "__Schema",
				"description": `A GraphQL Schema defines the capabilities of a GraphQL ` +
					`server. It exposes all available types and directives on ` +
					`the server, as well as the entry points for query, mutation, ` +
					`and subscription operations.`,
				"fields": []any{
					map[string]any{
						"name":        "types",
						"description": "A list of all types supported by this server.",
					},
					map[string]any{
						"name":        "queryType",
						"description": "The type that query operations will be rooted at.",
					},
					map[string]any{
						"name": "mutationType",
						"description": "If this server supports mutation, the type that " +
							"mutation operations will be rooted at.",
					},
					map[string]any{
						"name": "subscriptionType",
						"description": "If this server supports subscription, the type that " +
							"subscription operations will be rooted at.",
					},
					map[string]any{
						"name":        "directives",
						"description": "A list of all directives supported by this server.",
					},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expected.Data.(map[string]any)) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestIntrospection_ExposesDescriptionsOnEnums(t *testing.T) {

	queryRoot := graphql.NewObject(graphql.ObjectConfig{
		Name: "QueryRoot",
		Fields: graphql.Fields{
			"onlyField": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: queryRoot,
	})
	if err != nil {
		t.Fatalf("Error creating Schema: %v", err.Error())
	}
	query := `
      {
        typeKindType: __type(name: "__TypeKind") {
          name,
          description,
          enumValues {
            name,
            description
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]any{
			"typeKindType": map[string]any{
				"name":        "__TypeKind",
				"description": "An enum describing what kind of type a given `__Type` is",
				"enumValues": []any{
					map[string]any{
						"name":        "SCALAR",
						"description": "Indicates this type is a scalar.",
					},
					map[string]any{
						"name":        "OBJECT",
						"description": "Indicates this type is an object. `fields` and `interfaces` are valid fields.",
					},
					map[string]any{
						"name":        "INTERFACE",
						"description": "Indicates this type is an interface. `fields` and `possibleTypes` are valid fields.",
					},
					map[string]any{
						"name":        "UNION",
						"description": "Indicates this type is a union. `possibleTypes` is a valid field.",
					},
					map[string]any{
						"name":        "ENUM",
						"description": "Indicates this type is an enum. `enumValues` is a valid field.",
					},
					map[string]any{
						"name":        "INPUT_OBJECT",
						"description": "Indicates this type is an input object. `inputFields` is a valid field.",
					},
					map[string]any{
						"name":        "LIST",
						"description": "Indicates this type is a list. `ofType` is a valid field.",
					},
					map[string]any{
						"name":        "NON_NULL",
						"description": "Indicates this type is a non-null. `ofType` is a valid field.",
					},
				},
			},
		},
	}
	result := g(t, graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if !testutil.ContainSubset(result.Data.(map[string]any), expected.Data.(map[string]any)) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
