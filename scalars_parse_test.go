package graphql_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
)

func TestTypeSystem_Scalar_ParseValueOutputDateTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2017-07-23T03:46:56.647Z")
	tests := []dateTimeSerializationTest{
		{nil, nil, false},
		{"", nil, true},
		{(*string)(nil), nil, false},
		{"2017-07-23", nil, true},
		{"2017-07-23T03:46:56.647Z", t1, false},
	}
	for _, test := range tests {
		val, err := graphql.DateTime.ParseValue(test.Value)
		if err != nil && !test.Fails {
			t.Fatalf("failed DateTime.ParseValue(%v(%v)), expected: %v, got %v", reflect.TypeOf(test.Value), test.Value, test.Expected, err)
		}
		if err == nil && test.Fails {
			t.Fatalf("failed DateTime.ParseValue(%v(%v)), should have failed", reflect.TypeOf(test.Value), test.Value)
		} else if val != test.Expected {
			reflectedValue := reflect.ValueOf(test.Value)
			t.Fatalf("failed DateTime.ParseValue(%v(%v)), expected: %v, got %v", reflectedValue.Type(), test.Value, test.Expected, val)
		}
	}
}

func TestTypeSystem_Scalar_ParseLiteralOutputDateTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2017-07-23T03:46:56.647Z")
	for name, testCase := range map[string]struct {
		Literal  ast.Value
		Expected any
		Fails    bool
	}{
		"String": {
			Literal: &ast.StringValue{
				Value: "2017-07-23T03:46:56.647Z",
			},
			Expected: t1,
		},
		"NotAString": {
			Literal:  &ast.IntValue{},
			Expected: nil,
			Fails:    true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			parsed, err := graphql.DateTime.ParseLiteral(testCase.Literal)
			if err != nil && !testCase.Fails {
				t.Fatalf("failed DateTime.ParseLiteral(%T(%v)), expected: %v, got %v", testCase.Literal, testCase.Literal, testCase.Expected, err)
			}
			if err == nil && testCase.Fails {
				t.Fatalf("failed DateTime.ParseLiteral(%T(%v)), should have failed", testCase.Literal, testCase.Literal)
			} else if parsed != testCase.Expected {
				t.Fatalf("failed DateTime.ParseLiteral(%T(%v)), expected: %v, got %v", testCase.Literal, testCase.Literal, parsed, testCase.Expected)
			}
		})
	}
}
