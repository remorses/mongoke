package fields

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
	goke "github.com/remorses/goke/src"
)

type applyGuardsOnDocumentParams struct {
	document  goke.Map
	guards    []goke.AuthGuard
	jwt       jwt.MapClaims
	operation string
}

func applyGuardsOnDocument(p applyGuardsOnDocumentParams) (goke.Map, error) {
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

func hideFieldsFromDocument(documentMap goke.Map, toHide []string) goke.Map {
	if documentMap == nil {
		return nil
	}

	// clone the map
	copy := make(goke.Map, len(documentMap))
	for k, v := range documentMap {
		copy[k] = v
	}
	// remove the names from the copy
	for _, name := range toHide {
		_, ok := documentMap[name]
		if ok {
			delete(copy, name)
		}
	}
	return copy
}

func evaluateAuthPermission(guards []goke.AuthGuard, jwt jwt.MapClaims, document interface{}) (goke.AuthGuard, error) {
	// TODO if user is admin return the all permissions AuthGuard here
	// if guards are empty default to read permission
	if len(guards) == 0 {
		return goke.AuthGuard{
			AllowedOperations: []string{goke.Operations.READ},
		}, nil
	}
	for _, guard := range guards {
		res, err := guard.Evaluate(
			goke.Map{
				"jwt":      jwt,
				"document": document,
				// TODO more auth evaluation params like x, utility functions, ...
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
					goke.Operations.CREATE,
					goke.Operations.DELETE,
					goke.Operations.READ,
					goke.Operations.UPDATE,
				}
			}
			return guard, nil
		}
	}
	// default last permission is nothing when permissions list not empty and error
	permission := goke.AuthGuard{
		AllowedOperations: nil,
	}
	return permission, errors.New("no required permission for this resource")
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
