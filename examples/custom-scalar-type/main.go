package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
)

type CustomID struct {
	value string
}

func (id *CustomID) String() string {
	return id.value
}

func NewCustomID(v string) *CustomID {
	return &CustomID{value: v}
}

var CustomScalarType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "CustomScalarType",
	Description: "The `CustomScalarType` scalar type represents an ID Object.",
	// Serialize serializes `CustomID` to string.
	Serialize: func(value any) (any, error) {
		switch value := value.(type) {
		case CustomID:
			return value.String(), nil
		case *CustomID:
			v := *value
			return v.String(), nil
		default:
			return nil, fmt.Errorf("value must be of type CustomID, have %+v", value)
		}
	},
	// ParseValue parses GraphQL variables from `string` to `CustomID`.
	ParseValue: func(value any) (any, error) {
		switch value := value.(type) {
		case string:
			return NewCustomID(value), nil
		case *string:
			return NewCustomID(*value), nil
		default:
			return nil, fmt.Errorf("value must be of type string, have %+v", value)
		}
	},
	// ParseLiteral parses GraphQL AST value to `CustomID`.
	ParseLiteral: func(valueAST ast.Value) (any, error) {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return NewCustomID(valueAST.Value), nil
		default:
			return nil, fmt.Errorf("valueAST must be of type ast.StringValue, have %+v", valueAST)
		}
	},
})

type Customer struct {
	ID *CustomID `json:"id"`
}

var CustomerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Customer",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: CustomScalarType,
		},
	},
})

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"customers": &graphql.Field{
					Type: graphql.NewList(CustomerType),
					Args: graphql.FieldConfigArgument{
						&graphql.ArgumentConfig{
							Name: "id",
							Type: CustomScalarType,
						},
					},
					Resolve: func(p graphql.ResolveParams) (any, error) {
						// id := p.Args["id"]
						// log.Printf("id from arguments: %+v", id)
						customers := []Customer{
							{ID: NewCustomID("fb278f2a4a13f")},
						}
						return customers, nil
					},
				},
			},
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	query := `
		query {
			customers {
				id
			}
		}
	`
	/*
		queryWithVariable := `
			query($id: CustomScalarType) {
				customers(id: $id) {
					id
				}
			}
		`
	*/
	/*
		queryWithArgument := `
			query {
				customers(id: "5b42ba57289") {
					id
				}
			}
		`
	*/
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		VariableValues: map[string]any{
			"id": "5b42ba57289",
		},
	})
	if len(result.Errors) > 0 {
		log.Fatal(result)
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
