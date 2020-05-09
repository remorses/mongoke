package mongoke

import (
	"context"
	"errors"

	"github.com/PaesslerAG/gval"
	jwt "github.com/dgrijalva/jwt-go"
)

func (guard AuthGuard) Evaluate(params Map) (interface{}, error) {
	if guard.eval == nil {
		eval, err := gval.Full().NewEvaluable(guard.Expression)
		if err != nil {
			return nil, err
		}
		guard.eval = eval
	}
	res, err := guard.eval(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type applyGuardsOnDocumentParams struct {
	document  interface{}
	guards    []AuthGuard
	jwt       jwt.MapClaims
	operation string
}

func applyGuardsOnDocument(p applyGuardsOnDocumentParams) (interface{}, error) {
	if p.document == nil {
		return nil, nil
	}

	guard, err := evaluateAuthPermission(p.guards, p.jwt, p.document)
	if err != nil {
		return nil, err
	}
	if !contains(guard.AllowedOperations, p.operation) {
		return nil, errors.New("cannot execute " + p.operation + " operation with current user permissions")
	}
	if len(guard.HideFields) != 0 {
		p.document = hideFieldsFromDocument(p.document, guard.HideFields)
	}
	return p.document, nil
}

func hideFieldsFromDocument(document interface{}, toHide []string) interface{} {
	if document == nil {
		return nil
	}
	documentMap, ok := document.(Map)
	if !ok {
		return document
	}
	for _, name := range toHide {
		_, ok := documentMap[name]
		if ok {
			delete(documentMap, name)
		}
	}
	return documentMap

}

func evaluateAuthPermission(guards []AuthGuard, jwt jwt.MapClaims, document interface{}) (AuthGuard, error) {
	// TODO if user if admin return the max AuthGuard here
	// if guards are empty default to read permission
	if len(guards) == 0 {
		return AuthGuard{
			AllowedOperations: []string{Operations.READ},
		}, nil
	}
	for _, guard := range guards {
		res, err := guard.Evaluate(
			Map{
				"jwt":      jwt,
				"document": document,
				// TODO more evaluation params like x, utility functions, ...
			},
		)
		if err != nil {
			println("error while evaluating expression " + guard.Expression)
			continue
		}
		if res == true {
			// default allowed operations is every operation
			if len(guard.AllowedOperations) == 0 {
				guard.AllowedOperations = []string{
					Operations.CREATE,
					Operations.DELETE,
					Operations.READ,
					Operations.UPDATE,
				}
			}
			return guard, nil
		}
	}
	// default last permission is nothing when permissions list not empty
	permission := AuthGuard{
		AllowedOperations: nil,
	}
	return permission, nil // TODO should return an error here?
}
