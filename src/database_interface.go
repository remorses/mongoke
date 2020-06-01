package mongoke

import "context"

//go:generate moq -pkg mock -out mock/database_interface_mock.go . DatabaseInterface
type DatabaseInterface interface {
	// FindOne(p FindOneParams) (interface{}, error)
	// FindMany should return p.First + 1 nodes, or p.Last + 1 nodes, so mongoke can compute `hasNextPage` and `hasPreviousPage`
	FindMany(ctx context.Context, p FindManyParams) ([]Map, error)
	// TODO add InsertMany
	// TODO add UpdateMany
	// TODO add UpdateOne
}

// type FindOneParams struct {
// 	Collection  string
// 	DatabaseUri string
// 	Where       map[string]Filter `mapstructure:"where"`
// }

type FindManyParams struct {
	Collection  string
	DatabaseUri string
	Where       map[string]Filter `mapstructure:"where"`
	Limit       int               `mapstructure:"limit"`
	Offset      int               `mapstructure:"offset"`
	OrderBy     map[string]int    `mapstructure:"orderBy"`
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
