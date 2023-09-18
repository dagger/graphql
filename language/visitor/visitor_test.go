package visitor_test

import (
	"io/ioutil"
	"reflect"
	"testing"

	"fmt"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
	"github.com/dagger/graphql/language/kinds"
	"github.com/dagger/graphql/language/parser"
	"github.com/dagger/graphql/language/printer"
	"github.com/dagger/graphql/language/visitor"
	"github.com/dagger/graphql/testutil"
)

func parse(t *testing.T, query string) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			NoLocation: true,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}

func TestVisitor_AllowsEditingANodeBothOnEnterAndOnLeave(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	var selectionSet *ast.SelectionSet

	expectedQuery := `{ a, b, c { a, b, c } }`
	expectedAST := parse(t, expectedQuery)

	visited := map[string]bool{
		"didEnter": false,
		"didLeave": false,
	}

	expectedVisited := map[string]bool{
		"didEnter": true,
		"didLeave": true,
	}

	v := &visitor.VisitorOptions{

		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			kinds.OperationDefinition: {
				Enter: func(p visitor.VisitFuncParams) (string, any) {
					if node, ok := p.Node.(*ast.OperationDefinition); ok {
						selectionSet = node.SelectionSet
						visited["didEnter"] = true
						return visitor.ActionUpdate, ast.NewOperationDefinition(&ast.OperationDefinition{
							Loc:                 node.Loc,
							Operation:           node.Operation,
							Name:                node.Name,
							VariableDefinitions: node.VariableDefinitions,
							Directives:          node.Directives,
							SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
								Selections: []ast.Selection{},
							}),
						})
					}
					return visitor.ActionNoChange, nil
				},
				Leave: func(p visitor.VisitFuncParams) (string, any) {
					if node, ok := p.Node.(*ast.OperationDefinition); ok {
						visited["didLeave"] = true
						return visitor.ActionUpdate, ast.NewOperationDefinition(&ast.OperationDefinition{
							Loc:                 node.Loc,
							Operation:           node.Operation,
							Name:                node.Name,
							VariableDefinitions: node.VariableDefinitions,
							Directives:          node.Directives,
							SelectionSet:        selectionSet,
						})
					}
					return visitor.ActionNoChange, nil
				},
			},
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(expectedAST, editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}

}
func TestVisitor_AllowsEditingTheRootNodeOnEnterAndOnLeave(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	definitions := astDoc.Definitions

	expectedQuery := `{ a, b, c { a, b, c } }`
	expectedAST := parse(t, expectedQuery)

	visited := map[string]bool{
		"didEnter": false,
		"didLeave": false,
	}

	expectedVisited := map[string]bool{
		"didEnter": true,
		"didLeave": true,
	}

	v := &visitor.VisitorOptions{

		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			kinds.Document: {
				Enter: func(p visitor.VisitFuncParams) (string, any) {
					if node, ok := p.Node.(*ast.Document); ok {
						visited["didEnter"] = true
						return visitor.ActionUpdate, ast.NewDocument(&ast.Document{
							Loc:         node.Loc,
							Definitions: []ast.Node{},
						})
					}
					return visitor.ActionNoChange, nil
				},
				Leave: func(p visitor.VisitFuncParams) (string, any) {
					if node, ok := p.Node.(*ast.Document); ok {
						visited["didLeave"] = true
						return visitor.ActionUpdate, ast.NewDocument(&ast.Document{
							Loc:         node.Loc,
							Definitions: definitions,
						})
					}
					return visitor.ActionNoChange, nil
				},
			},
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(expectedAST, editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}
func TestVisitor_AllowsForEditingOnEnter(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedQuery := `{ a,    c { a,    c } }`
	expectedAST := parse(t, expectedQuery)
	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(expectedAST, editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}

}
func TestVisitor_AllowsForEditingOnLeave(t *testing.T) {

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedQuery := `{ a,    c { a,    c } }`
	expectedAST := parse(t, expectedQuery)
	v := &visitor.VisitorOptions{
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAst := visitor.Visit(astDoc, v, nil)
	if !reflect.DeepEqual(expectedAST, editedAst) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedAST, editedAst))
	}
}

func TestVisitor_VisitsEditedNode(t *testing.T) {

	query := `{ a { x } }`
	astDoc := parse(t, query)

	addedField := &ast.Field{
		Kind: "Field",
		Name: &ast.Name{
			Kind:  "Name",
			Value: "__typename",
		},
	}

	didVisitAddedField := false
	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node.Name != nil && node.Name.Value == "a" {
					s := node.SelectionSet.Selections
					s = append(s, addedField)
					ss := node.SelectionSet
					ss.Selections = s
					return visitor.ActionUpdate, ast.NewField(&ast.Field{
						Kind:         "Field",
						SelectionSet: ss,
					})
				}
				if reflect.DeepEqual(node, addedField) {
					didVisitAddedField = true
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)
	if didVisitAddedField == false {
		t.Fatalf("Unexpected result, expected didVisitAddedField == true")
	}
}
func TestVisitor_AllowsSkippingASubTree(t *testing.T) {

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	visited := []any{}
	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "Name", "c"},
		[]any{"leave", "Field", nil},
		[]any{"leave", "SelectionSet", nil},
		[]any{"leave", "OperationDefinition", nil},
		[]any{"leave", "Document", nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
			case *ast.Field:
				visited = append(visited, []any{"enter", node.Kind, nil})
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionSkip, nil
				}
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_AllowsEarlyExitWhileVisiting(t *testing.T) {

	visited := []any{}

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "b"},
		[]any{"leave", "Name", "b"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "x"},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
				if node.Value == "x" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_AllowsEarlyExitWhileLeaving(t *testing.T) {

	visited := []any{}

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "b"},
		[]any{"leave", "Name", "b"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "x"},
		[]any{"leave", "Name", "x"},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
				if node.Value == "x" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_AllowsANamedFunctionsVisitorAPI(t *testing.T) {

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	visited := []any{}
	expectedVisited := []any{
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Name", "a"},
		[]any{"enter", "Name", "b"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Name", "x"},
		[]any{"leave", "SelectionSet", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "SelectionSet", nil},
	}

	v := &visitor.VisitorOptions{
		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			"Name": {
				Kind: func(p visitor.VisitFuncParams) (string, any) {
					switch node := p.Node.(type) {
					case *ast.Name:
						visited = append(visited, []any{"enter", node.Kind, node.Value})
					}
					return visitor.ActionNoChange, nil
				},
			},
			"SelectionSet": {
				Enter: func(p visitor.VisitFuncParams) (string, any) {
					switch node := p.Node.(type) {
					case *ast.SelectionSet:
						visited = append(visited, []any{"enter", node.Kind, nil})
					}
					return visitor.ActionNoChange, nil
				},
				Leave: func(p visitor.VisitFuncParams) (string, any) {
					switch node := p.Node.(type) {
					case *ast.SelectionSet:
						visited = append(visited, []any{"leave", node.Kind, nil})
					}
					return visitor.ActionNoChange, nil
				},
			},
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}
func TestVisitor_VisitsKitchenSink(t *testing.T) {
	b, err := ioutil.ReadFile("../../kitchen-sink.graphql")
	if err != nil {
		t.Fatalf("unable to load kitchen-sink.graphql")
	}

	query := string(b)
	astDoc := parse(t, query)

	visited := []any{}
	expectedVisited := []any{
		[]any{"enter", "Document", nil, nil},
		[]any{"enter", "OperationDefinition", 0, nil},
		[]any{"enter", "Name", "Name", "OperationDefinition"},
		[]any{"leave", "Name", "Name", "OperationDefinition"},
		[]any{"enter", "VariableDefinition", 0, nil},
		[]any{"enter", "Variable", "Variable", "VariableDefinition"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Variable", "VariableDefinition"},
		[]any{"enter", "Named", "Type", "VariableDefinition"},
		[]any{"enter", "Name", "Name", "Named"},
		[]any{"leave", "Name", "Name", "Named"},
		[]any{"leave", "Named", "Type", "VariableDefinition"},
		[]any{"leave", "VariableDefinition", 0, nil},
		[]any{"enter", "VariableDefinition", 1, nil},
		[]any{"enter", "Variable", "Variable", "VariableDefinition"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Variable", "VariableDefinition"},
		[]any{"enter", "Named", "Type", "VariableDefinition"},
		[]any{"enter", "Name", "Name", "Named"},
		[]any{"leave", "Name", "Name", "Named"},
		[]any{"leave", "Named", "Type", "VariableDefinition"},
		[]any{"enter", "EnumValue", "DefaultValue", "VariableDefinition"},
		[]any{"leave", "EnumValue", "DefaultValue", "VariableDefinition"},
		[]any{"leave", "VariableDefinition", 1, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Alias", "Field"},
		[]any{"leave", "Name", "Alias", "Field"},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "ListValue", "Value", "Argument"},
		[]any{"enter", "IntValue", 0, nil},
		[]any{"leave", "IntValue", 0, nil},
		[]any{"enter", "IntValue", 1, nil},
		[]any{"leave", "IntValue", 1, nil},
		[]any{"leave", "ListValue", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"enter", "InlineFragment", 1, nil},
		[]any{"enter", "Named", "TypeCondition", "InlineFragment"},
		[]any{"enter", "Name", "Name", "Named"},
		[]any{"leave", "Name", "Name", "Named"},
		[]any{"leave", "Named", "TypeCondition", "InlineFragment"},
		[]any{"enter", "Directive", 0, nil},
		[]any{"enter", "Name", "Name", "Directive"},
		[]any{"leave", "Name", "Name", "Directive"},
		[]any{"leave", "Directive", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "InlineFragment"},

		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"enter", "Field", 1, nil},
		[]any{"enter", "Name", "Alias", "Field"},
		[]any{"leave", "Name", "Alias", "Field"},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "IntValue", "Value", "Argument"},
		[]any{"leave", "IntValue", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"enter", "Argument", 1, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "Variable", "Value", "Argument"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Value", "Argument"},
		[]any{"leave", "Argument", 1, nil},
		[]any{"enter", "Directive", 0, nil},
		[]any{"enter", "Name", "Name", "Directive"},
		[]any{"leave", "Name", "Name", "Directive"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "Variable", "Value", "Argument"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"leave", "Directive", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"enter", "FragmentSpread", 1, nil},
		[]any{"enter", "Name", "Name", "FragmentSpread"},
		[]any{"leave", "Name", "Name", "FragmentSpread"},
		[]any{"leave", "FragmentSpread", 1, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 1, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]any{"leave", "InlineFragment", 1, nil},
		[]any{"enter", "InlineFragment", 2, nil},
		[]any{"enter", "Directive", 0, nil},
		[]any{"enter", "Name", "Name", "Directive"},
		[]any{"leave", "Name", "Name", "Directive"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "Variable", "Value", "Argument"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},

		[]any{"leave", "Variable", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"leave", "Directive", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]any{"leave", "InlineFragment", 2, nil},
		[]any{"enter", "InlineFragment", 3, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "InlineFragment"},
		[]any{"leave", "InlineFragment", 3, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"leave", "OperationDefinition", 0, nil},
		[]any{"enter", "OperationDefinition", 1, nil},
		[]any{"enter", "Name", "Name", "OperationDefinition"},
		[]any{"leave", "Name", "Name", "OperationDefinition"},
		[]any{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "IntValue", "Value", "Argument"},
		[]any{"leave", "IntValue", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"enter", "Directive", 0, nil},
		[]any{"enter", "Name", "Name", "Directive"},
		[]any{"leave", "Name", "Name", "Directive"},
		[]any{"leave", "Directive", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"leave", "OperationDefinition", 1, nil},
		[]any{"enter", "OperationDefinition", 2, nil},
		[]any{"enter", "Name", "Name", "OperationDefinition"},
		[]any{"leave", "Name", "Name", "OperationDefinition"},
		[]any{"enter", "VariableDefinition", 0, nil},
		[]any{"enter", "Variable", "Variable", "VariableDefinition"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},

		[]any{"leave", "Variable", "Variable", "VariableDefinition"},
		[]any{"enter", "Named", "Type", "VariableDefinition"},
		[]any{"enter", "Name", "Name", "Named"},
		[]any{"leave", "Name", "Name", "Named"},
		[]any{"leave", "Named", "Type", "VariableDefinition"},
		[]any{"leave", "VariableDefinition", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "Variable", "Value", "Argument"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"enter", "Field", 1, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "SelectionSet", "SelectionSet", "Field"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 1, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "Field"},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"leave", "OperationDefinition", 2, nil},
		[]any{"enter", "FragmentDefinition", 3, nil},
		[]any{"enter", "Name", "Name", "FragmentDefinition"},
		[]any{"leave", "Name", "Name", "FragmentDefinition"},
		[]any{"enter", "Named", "TypeCondition", "FragmentDefinition"},
		[]any{"enter", "Name", "Name", "Named"},
		[]any{"leave", "Name", "Name", "Named"},
		[]any{"leave", "Named", "TypeCondition", "FragmentDefinition"},
		[]any{"enter", "SelectionSet", "SelectionSet", "FragmentDefinition"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},

		[]any{"enter", "Variable", "Value", "Argument"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"enter", "Argument", 1, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "Variable", "Value", "Argument"},
		[]any{"enter", "Name", "Name", "Variable"},
		[]any{"leave", "Name", "Name", "Variable"},
		[]any{"leave", "Variable", "Value", "Argument"},
		[]any{"leave", "Argument", 1, nil},
		[]any{"enter", "Argument", 2, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "ObjectValue", "Value", "Argument"},
		[]any{"enter", "ObjectField", 0, nil},
		[]any{"enter", "Name", "Name", "ObjectField"},
		[]any{"leave", "Name", "Name", "ObjectField"},
		[]any{"enter", "StringValue", "Value", "ObjectField"},
		[]any{"leave", "StringValue", "Value", "ObjectField"},
		[]any{"leave", "ObjectField", 0, nil},
		[]any{"leave", "ObjectValue", "Value", "Argument"},
		[]any{"leave", "Argument", 2, nil},
		[]any{"leave", "Field", 0, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "FragmentDefinition"},
		[]any{"leave", "FragmentDefinition", 3, nil},
		[]any{"enter", "OperationDefinition", 4, nil},
		[]any{"enter", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"enter", "Field", 0, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"enter", "Argument", 0, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "BooleanValue", "Value", "Argument"},
		[]any{"leave", "BooleanValue", "Value", "Argument"},
		[]any{"leave", "Argument", 0, nil},
		[]any{"enter", "Argument", 1, nil},
		[]any{"enter", "Name", "Name", "Argument"},
		[]any{"leave", "Name", "Name", "Argument"},
		[]any{"enter", "BooleanValue", "Value", "Argument"},
		[]any{"leave", "BooleanValue", "Value", "Argument"},
		[]any{"leave", "Argument", 1, nil},
		[]any{"leave", "Field", 0, nil},
		[]any{"enter", "Field", 1, nil},
		[]any{"enter", "Name", "Name", "Field"},
		[]any{"leave", "Name", "Name", "Field"},
		[]any{"leave", "Field", 1, nil},
		[]any{"leave", "SelectionSet", "SelectionSet", "OperationDefinition"},
		[]any{"leave", "OperationDefinition", 4, nil},
		[]any{"leave", "Document", nil, nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case ast.Node:
				if p.Parent != nil {
					visited = append(visited, []any{"enter", node.GetKind(), p.Key, p.Parent.GetKind()})
				} else {
					visited = append(visited, []any{"enter", node.GetKind(), p.Key, nil})
				}
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case ast.Node:
				if p.Parent != nil {
					visited = append(visited, []any{"leave", node.GetKind(), p.Key, p.Parent.GetKind()})
				} else {
					visited = append(visited, []any{"leave", node.GetKind(), p.Key, nil})
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, v, nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsSkippingASubTree(t *testing.T) {

	// Note: nearly identical to the above test of the same test but
	// using visitInParallel.

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	visited := []any{}
	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "Name", "c"},
		[]any{"leave", "Field", nil},
		[]any{"leave", "SelectionSet", nil},
		[]any{"leave", "OperationDefinition", nil},
		[]any{"leave", "Document", nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
			case *ast.Field:
				visited = append(visited, []any{"enter", node.Kind, nil})
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionSkip, nil
				}
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsSkippingDifferentSubTrees(t *testing.T) {

	query := `{ a { x }, b { y} }`
	astDoc := parse(t, query)

	visited := []any{}
	expectedVisited := []any{
		[]any{"no-a", "enter", "Document", nil},
		[]any{"no-b", "enter", "Document", nil},
		[]any{"no-a", "enter", "OperationDefinition", nil},
		[]any{"no-b", "enter", "OperationDefinition", nil},
		[]any{"no-a", "enter", "SelectionSet", nil},
		[]any{"no-b", "enter", "SelectionSet", nil},
		[]any{"no-a", "enter", "Field", nil},
		[]any{"no-b", "enter", "Field", nil},
		[]any{"no-b", "enter", "Name", "a"},
		[]any{"no-b", "leave", "Name", "a"},
		[]any{"no-b", "enter", "SelectionSet", nil},
		[]any{"no-b", "enter", "Field", nil},
		[]any{"no-b", "enter", "Name", "x"},
		[]any{"no-b", "leave", "Name", "x"},
		[]any{"no-b", "leave", "Field", nil},
		[]any{"no-b", "leave", "SelectionSet", nil},
		[]any{"no-b", "leave", "Field", nil},
		[]any{"no-a", "enter", "Field", nil},
		[]any{"no-b", "enter", "Field", nil},
		[]any{"no-a", "enter", "Name", "b"},
		[]any{"no-a", "leave", "Name", "b"},
		[]any{"no-a", "enter", "SelectionSet", nil},
		[]any{"no-a", "enter", "Field", nil},
		[]any{"no-a", "enter", "Name", "y"},
		[]any{"no-a", "leave", "Name", "y"},
		[]any{"no-a", "leave", "Field", nil},
		[]any{"no-a", "leave", "SelectionSet", nil},
		[]any{"no-a", "leave", "Field", nil},
		[]any{"no-a", "leave", "SelectionSet", nil},
		[]any{"no-b", "leave", "SelectionSet", nil},
		[]any{"no-a", "leave", "OperationDefinition", nil},
		[]any{"no-b", "leave", "OperationDefinition", nil},
		[]any{"no-a", "leave", "Document", nil},
		[]any{"no-b", "leave", "Document", nil},
	}

	v := []*visitor.VisitorOptions{
		{
			Enter: func(p visitor.VisitFuncParams) (string, any) {
				switch node := p.Node.(type) {
				case *ast.Name:
					visited = append(visited, []any{"no-a", "enter", node.Kind, node.Value})
				case *ast.Field:
					visited = append(visited, []any{"no-a", "enter", node.Kind, nil})
					if node.Name != nil && node.Name.Value == "a" {
						return visitor.ActionSkip, nil
					}
				case ast.Node:
					visited = append(visited, []any{"no-a", "enter", node.GetKind(), nil})
				default:
					visited = append(visited, []any{"no-a", "enter", nil, nil})
				}
				return visitor.ActionNoChange, nil
			},
			Leave: func(p visitor.VisitFuncParams) (string, any) {
				switch node := p.Node.(type) {
				case *ast.Name:
					visited = append(visited, []any{"no-a", "leave", node.Kind, node.Value})
				case ast.Node:
					visited = append(visited, []any{"no-a", "leave", node.GetKind(), nil})
				default:
					visited = append(visited, []any{"no-a", "leave", nil, nil})
				}
				return visitor.ActionNoChange, nil
			},
		},
		{
			Enter: func(p visitor.VisitFuncParams) (string, any) {
				switch node := p.Node.(type) {
				case *ast.Name:
					visited = append(visited, []any{"no-b", "enter", node.Kind, node.Value})
				case *ast.Field:
					visited = append(visited, []any{"no-b", "enter", node.Kind, nil})
					if node.Name != nil && node.Name.Value == "b" {
						return visitor.ActionSkip, nil
					}
				case ast.Node:
					visited = append(visited, []any{"no-b", "enter", node.GetKind(), nil})
				default:
					visited = append(visited, []any{"no-b", "enter", nil, nil})
				}
				return visitor.ActionNoChange, nil
			},
			Leave: func(p visitor.VisitFuncParams) (string, any) {
				switch node := p.Node.(type) {
				case *ast.Name:
					visited = append(visited, []any{"no-b", "leave", node.Kind, node.Value})
				case ast.Node:
					visited = append(visited, []any{"no-b", "leave", node.GetKind(), nil})
				default:
					visited = append(visited, []any{"no-b", "leave", nil, nil})
				}
				return visitor.ActionNoChange, nil
			},
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v...), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsEarlyExitWhileVisiting(t *testing.T) {

	// Note: nearly identical to the above test of the same test but
	// using visitInParallel.

	visited := []any{}

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "b"},
		[]any{"leave", "Name", "b"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "x"},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
				if node.Value == "x" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsEarlyExitFromDifferentPoints(t *testing.T) {

	visited := []any{}

	query := `{ a { y }, b { x } }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"break-a", "enter", "Document", nil},
		[]any{"break-b", "enter", "Document", nil},
		[]any{"break-a", "enter", "OperationDefinition", nil},
		[]any{"break-b", "enter", "OperationDefinition", nil},
		[]any{"break-a", "enter", "SelectionSet", nil},
		[]any{"break-b", "enter", "SelectionSet", nil},
		[]any{"break-a", "enter", "Field", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-a", "enter", "Name", "a"},
		[]any{"break-b", "enter", "Name", "a"},
		[]any{"break-b", "leave", "Name", "a"},
		[]any{"break-b", "enter", "SelectionSet", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-b", "enter", "Name", "y"},
		[]any{"break-b", "leave", "Name", "y"},
		[]any{"break-b", "leave", "Field", nil},
		[]any{"break-b", "leave", "SelectionSet", nil},
		[]any{"break-b", "leave", "Field", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-b", "enter", "Name", "b"},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"break-a", "enter", node.Kind, node.Value})
				if node != nil && node.Value == "a" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []any{"break-a", "enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-a", "enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"break-a", "leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"break-a", "leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-a", "leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	v2 := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"break-b", "enter", node.Kind, node.Value})
				if node != nil && node.Value == "b" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []any{"break-b", "enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-b", "enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"break-b", "leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"break-b", "leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-b", "leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v, v2), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsEarlyExitWhileLeaving(t *testing.T) {

	visited := []any{}

	query := `{ a, b { x }, c }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "b"},
		[]any{"leave", "Name", "b"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "x"},
		[]any{"leave", "Name", "x"},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
				if node.Value == "x" {
					return visitor.ActionBreak, nil
				}
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsEarlyExitFromLeavingDifferentPoints(t *testing.T) {

	visited := []any{}

	query := `{ a { y }, b { x } }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"break-a", "enter", "Document", nil},
		[]any{"break-b", "enter", "Document", nil},
		[]any{"break-a", "enter", "OperationDefinition", nil},
		[]any{"break-b", "enter", "OperationDefinition", nil},
		[]any{"break-a", "enter", "SelectionSet", nil},
		[]any{"break-b", "enter", "SelectionSet", nil},
		[]any{"break-a", "enter", "Field", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-a", "enter", "Name", "a"},
		[]any{"break-b", "enter", "Name", "a"},
		[]any{"break-a", "leave", "Name", "a"},
		[]any{"break-b", "leave", "Name", "a"},
		[]any{"break-a", "enter", "SelectionSet", nil},
		[]any{"break-b", "enter", "SelectionSet", nil},
		[]any{"break-a", "enter", "Field", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-a", "enter", "Name", "y"},
		[]any{"break-b", "enter", "Name", "y"},
		[]any{"break-a", "leave", "Name", "y"},
		[]any{"break-b", "leave", "Name", "y"},
		[]any{"break-a", "leave", "Field", nil},
		[]any{"break-b", "leave", "Field", nil},
		[]any{"break-a", "leave", "SelectionSet", nil},
		[]any{"break-b", "leave", "SelectionSet", nil},
		[]any{"break-a", "leave", "Field", nil},
		[]any{"break-b", "leave", "Field", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-b", "enter", "Name", "b"},
		[]any{"break-b", "leave", "Name", "b"},
		[]any{"break-b", "enter", "SelectionSet", nil},
		[]any{"break-b", "enter", "Field", nil},
		[]any{"break-b", "enter", "Name", "x"},
		[]any{"break-b", "leave", "Name", "x"},
		[]any{"break-b", "leave", "Field", nil},
		[]any{"break-b", "leave", "SelectionSet", nil},
		[]any{"break-b", "leave", "Field", nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"break-a", "enter", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"break-a", "enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-a", "enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				visited = append(visited, []any{"break-a", "leave", node.GetKind(), nil})
				if node.Name != nil && node.Name.Value == "a" {
					return visitor.ActionBreak, nil
				}
			case *ast.Name:
				visited = append(visited, []any{"break-a", "leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"break-a", "leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-a", "leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	v2 := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"break-b", "enter", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"break-b", "enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-b", "enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				visited = append(visited, []any{"break-b", "leave", node.GetKind(), nil})
				if node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionBreak, nil
				}
			case *ast.Name:
				visited = append(visited, []any{"break-b", "leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"break-b", "leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"break-b", "leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v, v2), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsForEditingOnEnter(t *testing.T) {

	visited := []any{}

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "Name", "c"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "Name", "c"},
		[]any{"leave", "Field", nil},
		[]any{"leave", "SelectionSet", nil},
		[]any{"leave", "Field", nil},
		[]any{"leave", "SelectionSet", nil},
		[]any{"leave", "OperationDefinition", nil},
		[]any{"leave", "Document", nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node != nil && node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	v2 := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitInParallel(v, v2), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}
}

func TestVisitor_VisitInParallel_AllowsForEditingOnLeave(t *testing.T) {

	visited := []any{}

	query := `{ a, b, c { a, b, c } }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil},
		[]any{"enter", "OperationDefinition", nil},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "b"},
		[]any{"leave", "Name", "b"},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "Name", "c"},
		[]any{"enter", "SelectionSet", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "a"},
		[]any{"leave", "Name", "a"},
		[]any{"leave", "Field", nil},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "b"},
		[]any{"leave", "Name", "b"},
		[]any{"enter", "Field", nil},
		[]any{"enter", "Name", "c"},
		[]any{"leave", "Name", "c"},
		[]any{"leave", "Field", nil},
		[]any{"leave", "SelectionSet", nil},
		[]any{"leave", "Field", nil},
		[]any{"leave", "SelectionSet", nil},
		[]any{"leave", "OperationDefinition", nil},
		[]any{"leave", "Document", nil},
	}

	v := &visitor.VisitorOptions{
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Field:
				if node != nil && node.Name != nil && node.Name.Value == "b" {
					return visitor.ActionUpdate, nil
				}
			}
			return visitor.ActionNoChange, nil
		},
	}

	v2 := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"enter", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil})
			default:
				visited = append(visited, []any{"leave", nil, nil})
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAST := visitor.Visit(astDoc, visitor.VisitInParallel(v, v2), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}

	expectedEditedAST := parse(t, `{ a,    c { a,    c } }`)
	if !reflect.DeepEqual(editedAST, expectedEditedAST) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedEditedAST, editedAST))
	}
}

func TestVisitor_VisitWithTypeInfo_MaintainsTypeInfoDuringVisit(t *testing.T) {

	visited := []any{}

	typeInfo := graphql.NewTypeInfo(&graphql.TypeInfoConfig{
		Schema: testutil.TestSchema,
	})

	query := `{ human(id: 4) { name, pets { name }, unknown } }`
	astDoc := parse(t, query)

	expectedVisited := []any{
		[]any{"enter", "Document", nil, nil, nil, nil},
		[]any{"enter", "OperationDefinition", nil, nil, "QueryRoot", nil},
		[]any{"enter", "SelectionSet", nil, "QueryRoot", "QueryRoot", nil},
		[]any{"enter", "Field", nil, "QueryRoot", "Human", nil},
		[]any{"enter", "Name", "human", "QueryRoot", "Human", nil},
		[]any{"leave", "Name", "human", "QueryRoot", "Human", nil},
		[]any{"enter", "Argument", nil, "QueryRoot", "Human", "ID"},
		[]any{"enter", "Name", "id", "QueryRoot", "Human", "ID"},
		[]any{"leave", "Name", "id", "QueryRoot", "Human", "ID"},
		[]any{"enter", "IntValue", nil, "QueryRoot", "Human", "ID"},
		[]any{"leave", "IntValue", nil, "QueryRoot", "Human", "ID"},
		[]any{"leave", "Argument", nil, "QueryRoot", "Human", "ID"},
		[]any{"enter", "SelectionSet", nil, "Human", "Human", nil},
		[]any{"enter", "Field", nil, "Human", "String", nil},
		[]any{"enter", "Name", "name", "Human", "String", nil},
		[]any{"leave", "Name", "name", "Human", "String", nil},
		[]any{"leave", "Field", nil, "Human", "String", nil},
		[]any{"enter", "Field", nil, "Human", "[Pet]", nil},
		[]any{"enter", "Name", "pets", "Human", "[Pet]", nil},
		[]any{"leave", "Name", "pets", "Human", "[Pet]", nil},
		[]any{"enter", "SelectionSet", nil, "Pet", "[Pet]", nil},
		[]any{"enter", "Field", nil, "Pet", "String", nil},
		[]any{"enter", "Name", "name", "Pet", "String", nil},
		[]any{"leave", "Name", "name", "Pet", "String", nil},
		[]any{"leave", "Field", nil, "Pet", "String", nil},
		[]any{"leave", "SelectionSet", nil, "Pet", "[Pet]", nil},
		[]any{"leave", "Field", nil, "Human", "[Pet]", nil},
		[]any{"enter", "Field", nil, "Human", nil, nil},
		[]any{"enter", "Name", "unknown", "Human", nil, nil},
		[]any{"leave", "Name", "unknown", "Human", nil, nil},
		[]any{"leave", "Field", nil, "Human", nil, nil},
		[]any{"leave", "SelectionSet", nil, "Human", "Human", nil},
		[]any{"leave", "Field", nil, "QueryRoot", "Human", nil},
		[]any{"leave", "SelectionSet", nil, "QueryRoot", "QueryRoot", nil},
		[]any{"leave", "OperationDefinition", nil, nil, "QueryRoot", nil},
		[]any{"leave", "Document", nil, nil, nil, nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			var parentType any
			var ttype any
			var inputType any

			if typeInfo.ParentType() != nil {
				parentType = fmt.Sprintf("%v", typeInfo.ParentType())
			}
			if typeInfo.Type() != nil {
				ttype = fmt.Sprintf("%v", typeInfo.Type())
			}
			if typeInfo.InputType() != nil {
				inputType = fmt.Sprintf("%v", typeInfo.InputType())
			}

			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value, parentType, ttype, inputType})
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil, parentType, ttype, inputType})
			default:
				visited = append(visited, []any{"enter", nil, nil, parentType, ttype, inputType})
			}
			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			var parentType any
			var ttype any
			var inputType any

			if typeInfo.ParentType() != nil {
				parentType = fmt.Sprintf("%v", typeInfo.ParentType())
			}
			if typeInfo.Type() != nil {
				ttype = fmt.Sprintf("%v", typeInfo.Type())
			}
			if typeInfo.InputType() != nil {
				inputType = fmt.Sprintf("%v", typeInfo.InputType())
			}

			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value, parentType, ttype, inputType})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil, parentType, ttype, inputType})
			default:
				visited = append(visited, []any{"leave", nil, nil, parentType, ttype, inputType})
			}
			return visitor.ActionNoChange, nil
		},
	}

	_ = visitor.Visit(astDoc, visitor.VisitWithTypeInfo(typeInfo, v), nil)

	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}

}

func TestVisitor_VisitWithTypeInfo_MaintainsTypeInfoDuringEdit(t *testing.T) {

	visited := []any{}

	typeInfo := graphql.NewTypeInfo(&graphql.TypeInfoConfig{
		Schema: testutil.TestSchema,
	})

	astDoc := parse(t, `{ human(id: 4) { name, pets }, alien }`)

	expectedVisited := []any{
		[]any{"enter", "Document", nil, nil, nil, nil},
		[]any{"enter", "OperationDefinition", nil, nil, "QueryRoot", nil},
		[]any{"enter", "SelectionSet", nil, "QueryRoot", "QueryRoot", nil},
		[]any{"enter", "Field", nil, "QueryRoot", "Human", nil},
		[]any{"enter", "Name", "human", "QueryRoot", "Human", nil},
		[]any{"leave", "Name", "human", "QueryRoot", "Human", nil},
		[]any{"enter", "Argument", nil, "QueryRoot", "Human", "ID"},
		[]any{"enter", "Name", "id", "QueryRoot", "Human", "ID"},
		[]any{"leave", "Name", "id", "QueryRoot", "Human", "ID"},
		[]any{"enter", "IntValue", nil, "QueryRoot", "Human", "ID"},
		[]any{"leave", "IntValue", nil, "QueryRoot", "Human", "ID"},
		[]any{"leave", "Argument", nil, "QueryRoot", "Human", "ID"},
		[]any{"enter", "SelectionSet", nil, "Human", "Human", nil},
		[]any{"enter", "Field", nil, "Human", "String", nil},
		[]any{"enter", "Name", "name", "Human", "String", nil},
		[]any{"leave", "Name", "name", "Human", "String", nil},
		[]any{"leave", "Field", nil, "Human", "String", nil},
		[]any{"enter", "Field", nil, "Human", "[Pet]", nil},
		[]any{"enter", "Name", "pets", "Human", "[Pet]", nil},
		[]any{"leave", "Name", "pets", "Human", "[Pet]", nil},
		[]any{"enter", "SelectionSet", nil, "Pet", "[Pet]", nil},
		[]any{"enter", "Field", nil, "Pet", "String!", nil},
		[]any{"enter", "Name", "__typename", "Pet", "String!", nil},
		[]any{"leave", "Name", "__typename", "Pet", "String!", nil},
		[]any{"leave", "Field", nil, "Pet", "String!", nil},
		[]any{"leave", "SelectionSet", nil, "Pet", "[Pet]", nil},
		[]any{"leave", "Field", nil, "Human", "[Pet]", nil},
		[]any{"leave", "SelectionSet", nil, "Human", "Human", nil},
		[]any{"leave", "Field", nil, "QueryRoot", "Human", nil},
		[]any{"enter", "Field", nil, "QueryRoot", "Alien", nil},
		[]any{"enter", "Name", "alien", "QueryRoot", "Alien", nil},
		[]any{"leave", "Name", "alien", "QueryRoot", "Alien", nil},
		[]any{"enter", "SelectionSet", nil, "Alien", "Alien", nil},
		[]any{"enter", "Field", nil, "Alien", "String!", nil},
		[]any{"enter", "Name", "__typename", "Alien", "String!", nil},
		[]any{"leave", "Name", "__typename", "Alien", "String!", nil},
		[]any{"leave", "Field", nil, "Alien", "String!", nil},
		[]any{"leave", "SelectionSet", nil, "Alien", "Alien", nil},
		[]any{"leave", "Field", nil, "QueryRoot", "Alien", nil},
		[]any{"leave", "SelectionSet", nil, "QueryRoot", "QueryRoot", nil},
		[]any{"leave", "OperationDefinition", nil, nil, "QueryRoot", nil},
		[]any{"leave", "Document", nil, nil, nil, nil},
	}

	v := &visitor.VisitorOptions{
		Enter: func(p visitor.VisitFuncParams) (string, any) {
			var parentType any
			var ttype any
			var inputType any

			if typeInfo.ParentType() != nil {
				parentType = fmt.Sprintf("%v", typeInfo.ParentType())
			}
			if typeInfo.Type() != nil {
				ttype = fmt.Sprintf("%v", typeInfo.Type())
			}
			if typeInfo.InputType() != nil {
				inputType = fmt.Sprintf("%v", typeInfo.InputType())
			}

			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"enter", node.Kind, node.Value, parentType, ttype, inputType})
			case *ast.Field:
				visited = append(visited, []any{"enter", node.GetKind(), nil, parentType, ttype, inputType})

				// Make a query valid by adding missing selection sets.
				if node.SelectionSet == nil && graphql.IsCompositeType(graphql.GetNamed(typeInfo.Type())) {
					return visitor.ActionUpdate, ast.NewField(&ast.Field{
						Alias:      node.Alias,
						Name:       node.Name,
						Arguments:  node.Arguments,
						Directives: node.Directives,
						SelectionSet: ast.NewSelectionSet(&ast.SelectionSet{
							Selections: []ast.Selection{
								ast.NewField(&ast.Field{
									Name: ast.NewName(&ast.Name{
										Value: "__typename",
									}),
								}),
							},
						}),
					})
				}
			case ast.Node:
				visited = append(visited, []any{"enter", node.GetKind(), nil, parentType, ttype, inputType})
			default:
				visited = append(visited, []any{"enter", nil, nil, parentType, ttype, inputType})
			}

			return visitor.ActionNoChange, nil
		},
		Leave: func(p visitor.VisitFuncParams) (string, any) {
			var parentType any
			var ttype any
			var inputType any

			if typeInfo.ParentType() != nil {
				parentType = fmt.Sprintf("%v", typeInfo.ParentType())
			}
			if typeInfo.Type() != nil {
				ttype = fmt.Sprintf("%v", typeInfo.Type())
			}
			if typeInfo.InputType() != nil {
				inputType = fmt.Sprintf("%v", typeInfo.InputType())
			}

			switch node := p.Node.(type) {
			case *ast.Name:
				visited = append(visited, []any{"leave", node.Kind, node.Value, parentType, ttype, inputType})
			case ast.Node:
				visited = append(visited, []any{"leave", node.GetKind(), nil, parentType, ttype, inputType})
			default:
				visited = append(visited, []any{"leave", nil, nil, parentType, ttype, inputType})
			}
			return visitor.ActionNoChange, nil
		},
	}

	editedAST := visitor.Visit(astDoc, visitor.VisitWithTypeInfo(typeInfo, v), nil)

	editedASTQuery := printer.Print(editedAST.(ast.Node))
	expectedEditedASTQuery := printer.Print(parse(t, `{ human(id: 4) { name, pets { __typename } }, alien { __typename } }`))

	if !reflect.DeepEqual(editedASTQuery, expectedEditedASTQuery) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedEditedASTQuery, editedASTQuery))
	}
	if !reflect.DeepEqual(visited, expectedVisited) {
		t.Fatalf("Unexpected result, Diff: %v", testutil.Diff(expectedVisited, visited))
	}

}
