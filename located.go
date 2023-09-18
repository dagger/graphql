package graphql

import (
	"errors"

	"github.com/dagger/graphql/gqlerrors"
	"github.com/dagger/graphql/language/ast"
)

func NewLocatedError(err any, nodes []ast.Node) *gqlerrors.Error {
	return newLocatedError(err, nodes, nil)
}

func NewLocatedErrorWithPath(err any, nodes []ast.Node, path []any) *gqlerrors.Error {
	return newLocatedError(err, nodes, path)
}

func newLocatedError(err any, nodes []ast.Node, path []any) *gqlerrors.Error {
	if err, ok := err.(*gqlerrors.Error); ok {
		return err
	}

	var origError error
	message := "An unknown error occurred."
	if err, ok := err.(error); ok {
		message = err.Error()
		origError = err
	}
	if err, ok := err.(string); ok {
		message = err
		origError = errors.New(err)
	}
	stack := message
	return gqlerrors.NewErrorWithPath(
		message,
		nodes,
		stack,
		nil,
		[]int{},
		path,
		origError,
	)
}

func FieldASTsToNodeASTs(fieldASTs []*ast.Field) []ast.Node {
	nodes := []ast.Node{}
	for _, fieldAST := range fieldASTs {
		nodes = append(nodes, fieldAST)
	}
	return nodes
}
