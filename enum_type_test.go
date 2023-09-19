package graphql_test

import (
	"reflect"
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/location"
	"github.com/dagger/graphql/testutil"
)

var enumTypeTestColorType = graphql.NewEnum(graphql.EnumConfig{
	Name: "Color",
	Values: graphql.EnumValueConfigMap{
		"RED": &graphql.EnumValueConfig{
			Value: int64(0),
		},
		"GREEN": &graphql.EnumValueConfig{
			Value: int64(1),
		},
		"BLUE": &graphql.EnumValueConfig{
			Value: int64(2),
		},
	},
})
var enumTypeTestQueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"colorEnum": &graphql.Field{
			Type: enumTypeTestColorType,
			Args: graphql.FieldConfigArgument{
				&graphql.ArgumentConfig{
					Name: "fromEnum",
					Type: enumTypeTestColorType,
				},
				&graphql.ArgumentConfig{
					Name: "fromInt",
					Type: graphql.Int,
				},
				&graphql.ArgumentConfig{
					Name: "fromString",
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				if fromInt, ok := p.Args["fromInt"]; ok {
					return fromInt, nil
				}
				if fromString, ok := p.Args["fromString"]; ok {
					return fromString, nil
				}
				if fromEnum, ok := p.Args["fromEnum"]; ok {
					return fromEnum, nil
				}
				return nil, nil
			},
		},
		"colorInt": &graphql.Field{
			Type: graphql.Int,
			Args: graphql.FieldConfigArgument{
				&graphql.ArgumentConfig{
					Name: "fromEnum",
					Type: enumTypeTestColorType,
				},
				&graphql.ArgumentConfig{
					Name: "fromInt",
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				if fromInt, ok := p.Args["fromInt"]; ok {
					return fromInt, nil
				}
				if fromEnum, ok := p.Args["fromEnum"]; ok {
					return fromEnum, nil
				}
				return nil, nil
			},
		},
	},
})
var enumTypeTestMutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"favoriteEnum": &graphql.Field{
			Type: enumTypeTestColorType,
			Args: graphql.FieldConfigArgument{
				&graphql.ArgumentConfig{
					Name: "color",
					Type: enumTypeTestColorType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				if color, ok := p.Args["color"]; ok {
					return color, nil
				}
				return nil, nil
			},
		},
	},
})

var enumTypeTestSubscriptionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Subscription",
	Fields: graphql.Fields{
		"subscribeToEnum": &graphql.Field{
			Type: enumTypeTestColorType,
			Args: graphql.FieldConfigArgument{
				&graphql.ArgumentConfig{
					Name: "color",
					Type: enumTypeTestColorType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				if color, ok := p.Args["color"]; ok {
					return color, nil
				}
				return nil, nil
			},
		},
	},
})

var enumTypeTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:        enumTypeTestQueryType,
	Mutation:     enumTypeTestMutationType,
	Subscription: enumTypeTestSubscriptionType,
})

func executeEnumTypeTest(t *testing.T, query string) *graphql.Result {
	result := g(t, graphql.Params{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	return result
}
func executeEnumTypeTestWithParams(t *testing.T, query string, params map[string]any) *graphql.Result {
	result := g(t, graphql.Params{
		Schema:         enumTypeTestSchema,
		RequestString:  query,
		VariableValues: params,
	})
	return result
}
func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInput(t *testing.T) {
	query := "{ colorInt(fromEnum: GREEN) }"
	expected := &graphql.Result{
		Data: map[string]any{
			"colorInt": int64(1),
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_EnumMayBeOutputType(t *testing.T) {
	query := "{ colorEnum(fromInt: 1) }"
	expected := &graphql.Result{
		Data: map[string]any{
			"colorEnum": "GREEN",
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumMayBeBothInputAndOutputType(t *testing.T) {
	query := "{ colorEnum(fromEnum: GREEN) }"
	expected := &graphql.Result{
		Data: map[string]any{
			"colorEnum": "GREEN",
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringLiterals(t *testing.T) {
	query := `{ colorEnum(fromEnum: "GREEN") }`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Argument \"fromEnum\" has invalid value \"GREEN\".\n" +
					"Expected type \"Color\", found \"GREEN\".\n" +
					"Error: Enum Color cannot parse value: GREEN",
				Locations: []location.SourceLocation{
					{Line: 1, Column: 23},
				},
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptIncorrectInternalValue(t *testing.T) {
	query := `{ colorEnum(fromString: "GREEN") }`
	expected := &graphql.Result{
		Data: map[string]any{
			"colorEnum": nil,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueInPlaceOfEnumLiteral(t *testing.T) {
	query := `{ colorEnum(fromEnum: 1) }`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Argument \"fromEnum\" has invalid value 1.\n" +
					"Expected type \"Color\", found 1.\n" +
					"Error: Enum Color cannot parse value: 1",
				Locations: []location.SourceLocation{
					{Line: 1, Column: 23},
				},
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_DoesNotAcceptEnumLiteralInPlaceOfInt(t *testing.T) {
	query := `{ colorEnum(fromInt: GREEN) }`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Argument \"fromInt\" has invalid value GREEN.\n" +
					"Expected type \"Int\", found GREEN.\n" +
					"Error: cannot parse *ast.EnumValue to int",
				Locations: []location.SourceLocation{
					{Line: 1, Column: 23},
				},
			},
		},
	}
	result := executeEnumTypeTest(t, query)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsJSONStringAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]any{
		"color": "BLUE",
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"colorEnum": "BLUE",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInputArgumentsToMutations(t *testing.T) {
	query := `mutation x($color: Color!) { favoriteEnum(color: $color) }`
	params := map[string]any{
		"color": "GREEN",
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"favoriteEnum": "GREEN",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_AcceptsEnumLiteralsAsInputArgumentsToSubscriptions(t *testing.T) {
	query := `subscription x($color: Color!) { subscribeToEnum(color: $color) }`
	params := map[string]any{
		"color": "GREEN",
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"subscribeToEnum": "GREEN",
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueAsEnumVariable(t *testing.T) {
	query := `query test($color: Color!) { colorEnum(fromEnum: $color) }`
	params := map[string]any{
		"color": 2,
	}
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Variable \"$color\" got invalid value 2.\n" +
					"Expected type \"Color\", found \"2\".\n" +
					"Error: Enum Color cannot parse non string type: 2",
				Locations: []location.SourceLocation{
					{Line: 1, Column: 12},
				},
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptStringVariablesAsEnumInput(t *testing.T) {
	query := `query test($color: String!) { colorEnum(fromEnum: $color) }`
	params := map[string]any{
		"color": "BLUE",
	}
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Variable "$color" of type "String!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_DoesNotAcceptInternalValueVariableAsEnumInput(t *testing.T) {
	query := `query test($color: Int!) { colorEnum(fromEnum: $color) }`
	params := map[string]any{
		"color": 2,
	}
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Variable "$color" of type "Int!" used in position expecting type "Color".`,
			},
		},
	}
	result := executeEnumTypeTestWithParams(t, query, params)
	if !testutil.EqualErrorMessage(expected, result, 0) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayHaveAnInternalValueOfZero(t *testing.T) {
	query := `{
        colorEnum(fromEnum: RED)
        colorInt(fromEnum: RED)
      }`
	expected := &graphql.Result{
		Data: map[string]any{
			"colorEnum": "RED",
			"colorInt":  int64(0),
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestTypeSystem_EnumValues_EnumValueMayBeNullable(t *testing.T) {
	query := `{
        colorEnum
        colorInt
      }`
	expected := &graphql.Result{
		Data: map[string]any{
			"colorEnum": nil,
			"colorInt":  nil,
		},
	}
	result := executeEnumTypeTest(t, query)
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_EnumValueMayBePointer(t *testing.T) {
	var enumTypeTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"query": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "query",
						Fields: graphql.Fields{
							"color": &graphql.Field{
								Type: enumTypeTestColorType,
							},
							"foo": &graphql.Field{
								Description: "foo field",
								Type:        graphql.Int,
							},
						},
					}),
					Resolve: func(_ graphql.ResolveParams) (any, error) {
						one := int64(1)
						return struct {
							Color *int64 `graphql:"color"`
							Foo   *int64 `graphql:"foo"`
						}{&one, &one}, nil
					},
				},
			},
		}),
	})
	query := "{ query { color foo } }"
	expected := &graphql.Result{
		Data: map[string]any{
			"query": map[string]any{
				"color": "GREEN",
				"foo":   int64(1)}}}
	result := g(t, graphql.Params{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestTypeSystem_EnumValues_EnumValueMayBeNilPointer(t *testing.T) {
	var enumTypeTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"query": &graphql.Field{
					Type: graphql.NewObject(graphql.ObjectConfig{
						Name: "query",
						Fields: graphql.Fields{
							"color": &graphql.Field{
								Type: enumTypeTestColorType,
							},
						},
					}),
					Resolve: func(_ graphql.ResolveParams) (any, error) {
						return struct {
							Color *int `graphql:"color"`
						}{nil}, nil
					},
				},
			},
		}),
	})
	query := "{ query { color } }"
	expected := &graphql.Result{
		Data: map[string]any{
			"query": map[string]any{
				"color": nil,
			}},
	}
	result := g(t, graphql.Params{
		Schema:        enumTypeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
