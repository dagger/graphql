package graphql_test

import (
	"sort"
	"testing"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/location"
	"github.com/dagger/graphql/testutil"
)

var syncError = "sync"
var nonNullSyncError = "nonNullSync"
var promiseError = "promise"
var nonNullPromiseError = "nonNullPromise"

var throwingData = map[string]any{
	"sync": func() any {
		panic(syncError)
	},
	"nonNullSync": func() any {
		panic(nonNullSyncError)
	},
	"promise": func() any {
		panic(promiseError)
	},
	"nonNullPromise": func() any {
		panic(nonNullPromiseError)
	},
}

var nullingData = map[string]any{
	"sync": func() any {
		return nil
	},
	"nonNullSync": func() any {
		return nil
	},
	"promise": func() any {
		return nil
	},
	"nonNullPromise": func() any {
		return nil
	},
}

var dataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "DataType",
	Fields: graphql.Fields{
		"sync": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullSync": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"promise": &graphql.Field{
			Type: graphql.String,
		},
		"nonNullPromise": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var nonNullTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: dataType,
})

func init() {
	throwingData["nest"] = func() any {
		return throwingData
	}
	throwingData["nonNullNest"] = func() any {
		return throwingData
	}
	throwingData["promiseNest"] = func() any {
		return throwingData
	}
	throwingData["nonNullPromiseNest"] = func() any {
		return throwingData
	}

	nullingData["nest"] = func() any {
		return nullingData
	}
	nullingData["nonNullNest"] = func() any {
		return nullingData
	}
	nullingData["promiseNest"] = func() any {
		return nullingData
	}
	nullingData["nonNullPromiseNest"] = func() any {
		return nullingData
	}

	dataType.AddFieldConfig("nest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
	dataType.AddFieldConfig("promiseNest", &graphql.Field{
		Type: dataType,
	})
	dataType.AddFieldConfig("nonNullPromiseNest", &graphql.Field{
		Type: graphql.NewNonNull(dataType),
	})
}

// nulls a nullable field that panics
func TestNonNull_NullsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"sync": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: syncError,
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
				Path: []any{
					"sync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"promise": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{
						Line: 3, Column: 9,
					},
				},
				Path: []any{
					"promise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []any{
					"nest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []any{
					"nest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []any{
					"promiseNest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatThrowsInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{
						Line: 4, Column: 11,
					},
				},
				Path: []any{
					"promiseNest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAComplexTreeOfNullableFieldsThatThrow(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]any{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []any{
					"nest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 7, Column: 13},
				},
				Path: []any{
					"nest", "nest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 11, Column: 13},
				},
				Path: []any{
					"nest", "promiseNest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 16, Column: 11},
				},
				Path: []any{
					"promiseNest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 19, Column: 13},
				},
				Path: []any{
					"promiseNest", "nest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: syncError,
				Locations: []location.SourceLocation{
					{Line: 23, Column: 13},
				},
				Path: []any{
					"promiseNest", "promiseNest", "sync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 5, Column: 11},
				},
				Path: []any{
					"nest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 13},
				},
				Path: []any{
					"nest", "nest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 12, Column: 13},
				},
				Path: []any{
					"nest", "promiseNest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 17, Column: 11},
				},
				Path: []any{
					"promiseNest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 20, Column: 13},
				},
				Path: []any{
					"promiseNest", "nest", "promise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: promiseError,
				Locations: []location.SourceLocation{
					{Line: 24, Column: 13},
				},
				Path: []any{
					"promiseNest", "promiseNest", "promise",
				},
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldThrowsInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 19},
				},
				Path: []any{
					"nest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 19, Column: 19},
				},
				Path: []any{
					"promiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
				Path: []any{
					"anotherNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			}),
			gqlerrors.FormatError(gqlerrors.Error{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 41, Column: 19},
				},
				Path: []any{
					"anotherPromiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			}),
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Errors, result.Errors))
	}

}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNull(t *testing.T) {
	doc := `
      query Q {
        sync
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"sync": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsANullableFieldThatSynchronouslyReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promise
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"promise": nil,
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []any{
					"nest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsASynchronouslyReturnedObjectThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []any{
					"nest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullSynchronously(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullSync,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []any{
					"promiseNest",
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAnObjectReturnedInAPromiseThatContainsANonNullableFieldThatReturnsNullInAPromise(t *testing.T) {
	doc := `
      query Q {
        promiseNest {
          nonNullPromise,
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"promiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 4, Column: 11},
				},
				Path: []any{
					"promiseNest",
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsAComplexTreeOfNullableFieldsThatReturnNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
        promiseNest {
          sync
          promise
          nest {
            sync
            promise
          }
          promiseNest {
            sync
            promise
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest": map[string]any{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
			},
			"promiseNest": map[string]any{
				"sync":    nil,
				"promise": nil,
				"nest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
				"promiseNest": map[string]any{
					"sync":    nil,
					"promise": nil,
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected.Data, result.Data))
	}
}
func TestNonNull_NullsTheFirstNullableObjectAfterAFieldReturnsNullInALongChainOfFieldsThatAreNonNull(t *testing.T) {
	doc := `
      query Q {
        nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullSync
                }
              }
            }
          }
        }
        anotherNest: nest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
        anotherPromiseNest: promiseNest {
          nonNullNest {
            nonNullPromiseNest {
              nonNullNest {
                nonNullPromiseNest {
                  nonNullPromise
                }
              }
            }
          }
        }
      }
	`
	expected := &graphql.Result{
		Data: map[string]any{
			"nest":               nil,
			"promiseNest":        nil,
			"anotherNest":        nil,
			"anotherPromiseNest": nil,
		},
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 8, Column: 19},
				},
				Path: []any{
					"nest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			},
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 19, Column: 19},
				},
				Path: []any{
					"promiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullSync",
				},
			},
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 30, Column: 19},
				},
				Path: []any{
					"anotherNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			},
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 41, Column: 19},
				},
				Path: []any{
					"anotherPromiseNest", "nonNullNest", "nonNullPromiseNest", "nonNullNest",
					"nonNullPromiseNest", "nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	sort.Sort(gqlerrors.FormattedErrors(expected.Errors))
	sort.Sort(gqlerrors.FormattedErrors(result.Errors))
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}

func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldThrows(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullSyncError,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []any{
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldErrors(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: nonNullPromiseError,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []any{
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   throwingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldReturnsNull(t *testing.T) {
	doc := `
      query Q { nonNullSync }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullSync.`,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []any{
					"nonNullSync",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
func TestNonNull_NullsTheTopLevelIfSyncNonNullableFieldResolvesNull(t *testing.T) {
	doc := `
      query Q { nonNullPromise }
	`
	expected := &graphql.Result{
		Data: nil,
		Errors: []gqlerrors.FormattedError{
			{
				Message: `Cannot return null for non-nullable field DataType.nonNullPromise.`,
				Locations: []location.SourceLocation{
					{Line: 2, Column: 17},
				},
				Path: []any{
					"nonNullPromise",
				},
			},
		},
	}
	// parse query
	ast := testutil.TestParse(t, doc)

	// execute
	ep := graphql.ExecuteParams{
		Schema: nonNullTestSchema,
		AST:    ast,
		Root:   nullingData,
	}
	result := testutil.TestExecute(t, ep)
	if !testutil.EqualResults(expected, result) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, result))
	}
}
