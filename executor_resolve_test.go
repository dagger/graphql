package graphql_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/testutil"
)

func testSchema(t *testing.T, testField *graphql.Field) graphql.Schema {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"test": testField,
			},
		}),
	})
	if err != nil {
		t.Fatalf("Invalid schema: %v", err)
	}
	return schema
}

func TestExecutesResolveFunction_DefaultFunctionAccessesProperties(t *testing.T) {
	schema := testSchema(t, &graphql.Field{Type: graphql.String})

	source := map[string]any{
		"test": "testValue",
	}

	expected := map[string]any{
		"test": "testValue",
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test }`,
		RootObject:    source,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}
}

func TestExecutesResolveFunction_DefaultFunctionCallsMethods(t *testing.T) {
	schema := testSchema(t, &graphql.Field{Type: graphql.String})

	source := map[string]any{
		"test": func() any {
			return "testValue"
		},
	}

	expected := map[string]any{
		"test": "testValue",
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test }`,
		RootObject:    source,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}
}

func TestExecutesResolveFunction_UsesProvidedResolveFunction(t *testing.T) {
	schema := testSchema(t, &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			&graphql.ArgumentConfig{Name: "aStr", Type: graphql.String},
			&graphql.ArgumentConfig{Name: "aInt", Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			b, err := json.Marshal(p.Args)
			return string(b), err
		},
	})

	expected := map[string]any{
		"test": "{}",
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}

	expected = map[string]any{
		"test": `{"aStr":"String!"}`,
	}
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test(aStr: "String!") }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}

	expected = map[string]any{
		"test": `{"aInt":-123,"aStr":"String!"}`,
	}
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test(aInt: -123, aStr: "String!") }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}
}

func TestExecutesResolveFunction_UsesProvidedResolveFunction_SourceIsStruct_WithoutJSONTags(t *testing.T) {

	// For structs without JSON tags, it will map to upper-cased exported field names
	type SubObjectWithoutJSONTags struct {
		Str string
		Int int64
	}

	schema := testSchema(t, &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:        "SubObject",
			Description: "Maps GraphQL Object `SubObject` to Go struct `SubObjectWithoutJSONTags`",
			Fields: graphql.Fields{
				"Str": &graphql.Field{Type: graphql.String},
				"Int": &graphql.Field{Type: graphql.Int},
			},
		}),
		Args: graphql.FieldConfigArgument{
			&graphql.ArgumentConfig{Name: "aStr", Type: graphql.String},
			&graphql.ArgumentConfig{Name: "aInt", Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			aStr, _ := p.Args["aStr"].(string)
			aInt, _ := p.Args["aInt"].(int64)
			return &SubObjectWithoutJSONTags{
				Str: aStr,
				Int: aInt,
			}, nil
		},
	})

	expected := map[string]any{
		"test": map[string]any{
			"Str": "",
			"Int": int64(0),
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test { Str, Int } }`,
	})

	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}

	expected = map[string]any{
		"test": map[string]any{
			"Str": "String!",
			"Int": int64(0),
		},
	}
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test(aStr: "String!") { Str, Int } }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}

	expected = map[string]any{
		"test": map[string]any{
			"Str": "String!",
			"Int": int64(-123),
		},
	}
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test(aInt: -123, aStr: "String!") { Str, Int } }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}
}

func TestExecutesResolveFunction_UsesProvidedResolveFunction_SourceIsStruct_WithJSONTags(t *testing.T) {

	// For structs without JSON tags, it will map to upper-cased exported field names
	type SubObjectWithJSONTags struct {
		OtherField string `json:""`
		Str        string `json:"str"`
		Int        int64  `json:"int"`
	}

	schema := testSchema(t, &graphql.Field{
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:        "SubObject",
			Description: "Maps GraphQL Object `SubObject` to Go struct `SubObjectWithJSONTags`",
			Fields: graphql.Fields{
				"str": &graphql.Field{Type: graphql.String},
				"int": &graphql.Field{Type: graphql.Int},
			},
		}),
		Args: graphql.FieldConfigArgument{
			&graphql.ArgumentConfig{Name: "aStr", Type: graphql.String},
			&graphql.ArgumentConfig{Name: "aInt", Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			aStr, _ := p.Args["aStr"].(string)
			aInt, _ := p.Args["aInt"].(int64)
			return &SubObjectWithJSONTags{
				Str: aStr,
				Int: aInt,
			}, nil
		},
	})

	expected := map[string]any{
		"test": map[string]any{
			"str": "",
			"int": int64(0),
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test { str, int } }`,
	})

	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}

	expected = map[string]any{
		"test": map[string]any{
			"str": "String!",
			"int": int64(0),
		},
	}
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test(aStr: "String!") { str, int } }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}

	expected = map[string]any{
		"test": map[string]any{
			"str": "String!",
			"int": int64(-123),
		},
	}
	result = graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: `{ test(aInt: -123, aStr: "String!") { str, int } }`,
	})
	if !reflect.DeepEqual(expected, result.Data) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result.Data))
	}
}
