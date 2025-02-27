package printer_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/dagger/graphql/language/ast"
	"github.com/dagger/graphql/language/printer"
	"github.com/dagger/graphql/testutil"
)

func TestSchemaPrinter_PrintsMinimalAST(t *testing.T) {
	astDoc := ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: ast.NewName(&ast.Name{
			Value: "foo",
		}),
	})
	results := printer.Print(astDoc)
	expected := "scalar foo"
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, results))
	}
}

func TestSchemaPrinter_DoesNotAlterAST(t *testing.T) {
	b, err := ioutil.ReadFile("../../schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	astDocBefore := testutil.ASTToJSON(t, astDoc)

	_ = printer.Print(astDoc)

	astDocAfter := testutil.ASTToJSON(t, astDoc)

	_ = testutil.ASTToJSON(t, astDoc)

	if !reflect.DeepEqual(astDocAfter, astDocBefore) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(astDocAfter, astDocBefore))
	}
}

func TestSchemaPrinter_PrintsKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("../../schema-kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)
	expected := `schema {
  query: QueryType
  mutation: MutationType
}

type Foo implements Bar & Baz {
  one: Type
  two(argument: InputType!): Type
  three(argument: InputType, other: String): Int
  four(argument: String = "string"): String
  five(argument: [String] = ["string", "string"]): String
  six(argument: InputType = {key: "value"}): Type
}

type AnnotatedObject @onObject(arg: "value") {
  annotatedField(arg: Type = "default" @onArg): Type @onField
}

interface Bar {
  one: Type
  four(argument: String = "string"): String
}

interface AnnotatedInterface @onInterface {
  annotatedField(arg: Type @onArg): Type @onField
}

union Feed = Story | Article | Advert

union AnnotatedUnion @onUnion = A | B

scalar CustomScalar

scalar AnnotatedScalar @onScalar

enum Site {
  DESKTOP
  MOBILE
}

enum AnnotatedEnum @onEnum {
  ANNOTATED_VALUE @onEnumValue
  OTHER_VALUE
}

input InputType {
  key: String!
  answer: Int = 42
}

input AnnotatedInput @onInputObjectType {
  annotatedField: Type @onField
}

extend type Foo {
  seven(argument: [String]): Type
}

extend type Foo @onType {}

type NoFields {}

directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT

directive @include(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
`
	results := printer.Print(astDoc)
	if !reflect.DeepEqual(expected, results) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, results))
	}
}

func TestSchemaPrinter_PrintsAllDescriptions(t *testing.T) {
	b, err := ioutil.ReadFile("../../schema-all-descriptions.graphql")
	if err != nil {
		t.Fatalf("unable to load schema-all-descriptions.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)
	expected := `"""single line scalar description"""
scalar ScalarSingleLine

"""
multi line

scalar description
"""
scalar ScalarMultiLine

"""single line object description"""
type ObjectSingleLine {
  no_description: ID
  
  """single line field description"""
  single_line(a: ID, b: ID, c: ID, d: ID): ID
  
  """
  multi line
  
  field description
  """
  multi_line(
    a: ID
    
    """single line argument description"""
    b: ID
    
    """
    multi line
    
    field description
    """
    c: ID
    d: ID
  ): ID
}

"""
multi line

object description
"""
type ObjectMultiLine {
  foo: ID
}

"""single line interface description"""
interface InterfaceSingleLine {
  no_description: ID
  
  """single line field description"""
  single_line(a: ID, b: ID, c: ID, d: ID): ID
  
  """
  multi line
  
  field description
  """
  multi_line(
    a: ID
    
    """single line argument description"""
    b: ID
    
    """
    multi line
    
    argument description
    """
    c: ID
    d: ID
  ): ID
}

"""
multi line

interface description
"""
interface InterfaceMultiLine {
  foo: ID
}

"""single line union description"""
union UnionSingleLine = String | Int | Float | ID

"""
multi line

union description
"""
union UnionSingleLine = String | Int | Float | ID

"""single line enum description"""
enum EnumSingleLine {
  no_description
  
  """single line enum description"""
  single_line
  
  """
  multi line
  
  enum description
  """
  multi_line
  again_no_description
}

"""
multi line

enum description
"""
enum EnumMultiLine {
  foo
}

"""single line input description"""
input InputSingleLine {
  a: ID
  
  """single line argument description"""
  b: ID
  
  """
  multi line
  
  argument description
  """
  c: ID
  d: ID
}

"""
multi line

input description
"""
input InputMultiLine {
  foo: ID
}

"""single line directive description"""
directive @DirectiveSingleLine(
  a: ID
  
  """single line argument description"""
  b: ID
  
  """
  multi line
  
  argument description
  """
  c: ID
  d: ID
) on SCALAR

"""
multi line

directive description
"""
directive @DirectiveMultiLine on SCALAR
`
	results := printer.Print(astDoc)
	if !reflect.DeepEqual(expected, results) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expected, results))
	}
}
