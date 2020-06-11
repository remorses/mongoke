package goke

import (
	"context"

	"github.com/mitchellh/mapstructure"
)

//go:generate moq -pkg mock -out mock/database_interface_mock.go . DatabaseInterface
type DatabaseInterface interface {
	// FindOne(p FindOneParams) (interface{}, error)
	// FindMany should return p.First + 1 nodes, or p.Last + 1 nodes, so goke can compute `hasNextPage` and `hasPreviousPage`
	FindMany(ctx context.Context, p FindManyParams, hook TransformDocument) ([]Map, error)
	InsertMany(ctx context.Context, p InsertManyParams, hook TransformDocument) (NodesMutationPayload, error)
	UpdateOne(ctx context.Context, p UpdateParams, hook TransformDocument) (NodeMutationPayload, error)
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
	Set        Map       `mapstructure:"set" bson:"$set"`
	Where      WhereTree // `mapstructure:"where"`
}

type DeleteManyParams struct {
	Collection string
	Where      WhereTree //  `mapstructure:"where"`
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

func MakeWhereTree(where map[string]interface{}, initialMatch map[string]Filter) (WhereTree, error) {
	// for every k, v use mapstructure to map to a filter
	// if k is or, and, use mapstructure to map to an array of filters
	if initialMatch == nil {
		initialMatch = make(map[string]Filter)
	}
	tree := WhereTree{
		Match: initialMatch,
	}
	if where == nil {
		return tree, nil
	}
	for k, v := range where {
		if k == "and" || k == "or" {
			for _, item := range v.([]interface{}) {
				w, err := MakeWhereTree(item.(map[string]interface{}), nil)
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
