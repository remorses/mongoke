package goke

import (
	"context"
	"reflect"

	"github.com/PaesslerAG/gval"
	"github.com/mitchellh/mapstructure"
)

//go:generate moq -pkg mock -out mock/database_interface_mock.go . DatabaseInterface
type DatabaseInterface interface {
	// FindOne(p FindOneParams) (interface{}, error)
	FindMany(ctx context.Context, p FindManyParams, hook TransformDocument) ([]Map, error)
	InsertMany(ctx context.Context, p InsertManyParams, hook TransformDocument) (NodesMutationPayload, error)
	UpdateMany(ctx context.Context, p UpdateParams, hook TransformDocument) (NodesMutationPayload, error)
	DeleteMany(ctx context.Context, p DeleteManyParams, hook TransformDocument) (NodesMutationPayload, error)
}

// remove non accessible fields from a document and returns nil if document is not accessible
type TransformDocument = func(document Map) (Map, error)

// TODO instead of using `and` append directly to the Match, it is faster
func FilterDocuments(xs []Map, filter TransformDocument) ([]Map, error) {
	var result []Map
	for _, x := range xs {
		if filter != nil {
			res, err := filter(x)
			if err != nil {
				return nil, err
			}
			if res == nil {
				continue
			}
			result = append(result, res)
		} else {
			if x == nil {
				continue
			}
			result = append(result, x)
		}

	}
	return result, nil
}

func MapsToInterfaces(nodes []Map) []interface{} {
	var data = make([]interface{}, len(nodes))
	for i, x := range nodes {
		data[i] = x
	}
	return data
}

func ExtendWhereMatch(where WhereTree, match map[string]Filter) WhereTree {
	// the where is implicitly copied
	if match == nil {
		return where
	}
	if where.Match == nil {
		where.Match = match
		return where
	}
	where.And = append(where.And, WhereTree{
		Match: match,
	})
	return where
}

// type FindOneParams struct {
// 	Collection  string
// 	DatabaseUri string
// 	Where       map[string]Filter `mapstructure:"where"`
// }

type NodeMutationPayload struct {
	Returning     interface{} `json:"returning"`
	AffectedCount int         `json:"affectedCount"`
}

type NodesMutationPayload struct {
	Returning     []Map `json:"returning"`
	AffectedCount int   `json:"affectedCount"`
}

type WhereTree struct {
	Match map[string]Filter
	And   []WhereTree
	Or    []WhereTree
}

type FindManyParams struct {
	Collection string
	Where      WhereTree      // `mapstructure:"where"`
	Limit      int            `mapstructure:"limit"`
	Offset     int            `mapstructure:"offset"`
	OrderBy    map[string]int `mapstructure:"orderBy"`
}

type InsertManyParams struct {
	Collection string
	Data       []Map `mapstructure:"data"`
}

type UpdateParams struct {
	Collection string
	Set        Map `mapstructure:"set" bson:"$set"`
	Limit      int
	Where      WhereTree // `mapstructure:"where"`
}

type DeleteManyParams struct {
	Collection string
	Where      WhereTree //  `mapstructure:"where"`
	Limit      int
}

type Pagination struct {
	First  int    `mapstructure:first`
	Last   int    `mapstructure:last`
	After  string `mapstructure:after`
	Before string `mapstructure:before`
}

type Filter struct {
	Eq  interface{}   `bson:"$eq,omitempty"`
	Neq interface{}   `bson:"$ne,omitempty"`
	In  []interface{} `bson:"$in,omitempty"`
	Nin []interface{} `bson:"$nin,omitempty"`
	Gt  interface{}   `bson:"$gt,omitempty"`
	Lt  interface{}   `bson:"$lt,omitempty"`
	Gte interface{}   `bson:"$gte,omitempty"`
	Lte interface{}   `bson:"$lte,omitempty"`
}

func (self Filter) Interpolate(scope Map) (Filter, error) {
	v := reflect.ValueOf(self)
	// TODO i am using a map to make code more general, optimize using struct fields
	result := Map{}
	typeOfv := v.Type()
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i).Interface()
		name := typeOfv.Field(i).Name
		switch value.(type) {
		case string:
			// TODO evaluate only if matches a certain regex, like {{ }}
			evaluated, err := evaluate(value.(string), scope)
			if err != nil {
				return Filter{}, err
			}
			result[name] = evaluated
		default:
			result[name] = value
		}
	}
	var filter Filter
	err := mapstructure.Decode(result, &filter)
	if err != nil {
		return Filter{}, err
	}
	return filter, nil
}

func evaluate(expression string, scope Map) (interface{}, error) {
	// TODO save the evaluable in object state to optimize
	eval, err := gval.Full().NewEvaluable(expression)
	if err != nil {
		return "", err
	}
	res, err := eval(context.Background(), scope)
	if err != nil {
		return "", err
	}
	return res, nil
}

func MakeWhereTree(where map[string]interface{}) (WhereTree, error) {
	// for every k, v use mapstructure to map to a filter
	// if k is or, and, use mapstructure to map to an array of filters
	tree := WhereTree{
		Match: make(map[string]Filter),
	}
	if where == nil {
		return tree, nil
	}
	for k, v := range where {
		if k == "and" || k == "or" {
			for _, item := range v.([]interface{}) {
				w, err := MakeWhereTree(item.(map[string]interface{}))
				if err != nil {
					return tree, err
				}
				if k == "and" {
					tree.And = append(tree.And, w)
				} else {
					tree.Or = append(tree.Or, w)
				}
			}
			continue
		}
		var match Filter
		err := mapstructure.Decode(v, &match)
		if err != nil {
			return tree, err
		}
		tree.Match[k] = match
	}
	return tree, nil
}
