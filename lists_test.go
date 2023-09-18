package graphql_test

import (
	"reflect"
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/location"
	"github.com/dagger/graphql/testutil"
)

func checkList(t *testing.T, testType graphql.Type, testData any, expected *graphql.Result) {
	// TODO: uncomment t.Helper when support for go1.8 is dropped.
	//t.Helper()
	data := map[string]any{
		"test": testData,
	}

	dataType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DataType",
		Fields: graphql.Fields{
			"test": &graphql.Field{
				Type: testType,
			},
		},
	})
	dataType.AddFieldConfig("nest", &graphql.Field{
		Type: dataType,
		Resolve: func(p graphql.ResolveParams) (any, error) {
			return data, nil
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: dataType,
	})
	if err != nil {
		t.Fatalf("Error in schema %v", err.Error())
	}

	// parse query
	ast := testutil.TestParse(t, `{ nest { test } }`)

	// execute
	ep := graphql.ExecuteParams{
		Schema: schema,
		AST:    ast,
		Root:   data,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}

}

// Describe [T] Array<T>
func TestLists_ListOfNullableObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := []any{
		1, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableObjects_ContainsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := []any{
		1, nil, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), nil, int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableObjects_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T] Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_ListOfNullableFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableFunc_ContainsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			int64(1), nil, int64(2),
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), nil, int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableFunc_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return nil
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T] Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_ListOfNullableArrayOfFuncContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_ListOfNullableArrayOfFuncContainsNulls(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return nil, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), nil, int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T]! Array<T>
func TestLists_NonNullListOfNullableObjectsContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))
	data := []any{
		1, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableObjectsContainsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))
	data := []any{
		1, nil, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), nil, int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableObjectsReturnsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
				},
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T]! Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NonNullListOfNullableFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableFunc_ContainsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, nil, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), nil, int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableFunc_ReturnsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return nil
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T]! Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NonNullListOfNullableArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNullableArrayOfFunc_ContainsNulls(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return nil, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), nil, int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!] Array<T>
func TestLists_NullableListOfNonNullObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))
	data := []any{
		1, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullObjects_ContainsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))
	data := []any{
		1, nil, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
					1,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullObjects_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T!] Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NullableListOfNonNullFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullFunc_ContainsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, nil, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
					1,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullFunc_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return nil
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!] Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NullableListOfNonNullArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NullableListOfNonNullArrayOfFunc_ContainsNulls(t *testing.T) {
	ttype := graphql.NewList(graphql.NewNonNull(graphql.Int))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error){...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return nil, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		/*
			// TODO: Because thunks are called after the result map has been assembled,
			// we are not able to traverse up the tree until we find a nullable type,
			// so in this case the entire data is nil. Will need some significant code
			// restructure to restore this.
			Data: map[string]interface{}{
				"nest": map[string]interface{}{
					"test": nil,
				},
			},
		*/
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
					1,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!]! Array<T>
func TestLists_NonNullListOfNonNullObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))
	data := []any{
		1, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullObjects_ContainsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))
	data := []any{
		1, nil, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
					1,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullObjects_ReturnsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
				},
			},
		},
	}
	checkList(t, ttype, nil, expected)
}

// Describe [T!]! Func()Array<T> // equivalent to Promise<Array<T>>
func TestLists_NonNullListOfNonNullFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullFunc_ContainsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return []any{
			1, nil, 2,
		}
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
					1,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullFunc_ReturnsNull(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a function that return values
	// Note that its uses the expected signature `func() interface{} {...}`
	data := func() any {
		return nil
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

// Describe [T!]! Array<Func()<T>> // equivalent to Array<Promise<T>>
func TestLists_NonNullListOfNonNullArrayOfFunc_ContainsValues(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}
func TestLists_NonNullListOfNonNullArrayOfFunc_ContainsNulls(t *testing.T) {
	ttype := graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.Int)))

	// `data` is a slice of functions that return values
	// Note that its uses the expected signature `func() (interface{}, error) {...}`
	data := []any{
		func() (any, error) {
			return 1, nil
		},
		func() (any, error) {
			return nil, nil
		},
		func() (any, error) {
			return 2, nil
		},
	}
	expected := &graphql.Result{
		/*
			// TODO: Because thunks are called after the result map has been assembled,
			// we are not able to traverse up the tree until we find a nullable type,
			// so in this case the entire data is nil. Will need some significant code
			// restructure to restore this.
			Data: map[string]interface{}{
				"nest": nil,
			},
		*/
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: "Cannot return null for non-nullable field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
					1,
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

func TestLists_UserErrorExpectIterableButDidNotGetOne(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := "Not an iterable"
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: "User Error: expected iterable, but did not find one for field DataType.test.",
				Locations: []location.SourceLocation{
					{
						Line:   1,
						Column: 10,
					},
				},
				Path: []any{
					"nest",
					"test",
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

func TestLists_ArrayOfNullableObjects_ContainsValues(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	data := [2]any{
		1, 2,
	}
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": []any{
					int64(1), int64(2),
				},
			},
		},
	}
	checkList(t, ttype, data, expected)
}

func TestLists_ValueMayBeNilPointer(t *testing.T) {
	var listTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"list": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
					Resolve: func(_ graphql.ResolveParams) (any, error) {
						return []int(nil), nil
					},
				},
			},
		}),
	})
	query := "{ list }"
	expected := &graphql.Result{
		Data: map[string]any{
			"list": []any{},
		},
	}
	result := g(t, graphql.Params{
		Schema:        listTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestLists_NullableListOfInt_ReturnsNull(t *testing.T) {
	ttype := graphql.NewList(graphql.Int)
	type dataType *[]int
	var data dataType
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"test": nil,
			},
		},
	}
	checkList(t, ttype, data, expected)
}
