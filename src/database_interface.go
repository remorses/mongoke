package mongoke

//go:generate moq -out database_interface_mock.go . DatabaseInterface

type Map map[string]interface{}

type DatabaseInterface interface {
	FindOne(p FindOneParams) (interface{}, error)
	// FindMany should return p.First + 1 nodes, or p.Last + 1 nodes, so mongoke can compute `hasNextPage` and `hasPreviousPage`
	FindMany(p FindManyParams) ([]Map, error)
	// TODO add mutations in databaseFunctions
}

type FindOneParams struct {
	Collection  string
	DatabaseUri string
	Where       map[string]Filter `mapstructure:"where"`
}

type FindManyParams struct {
	Collection  string
	DatabaseUri string
	Where       map[string]Filter `mapstructure:"where"`
	Pagination  Pagination
	CursorField string `mapstructure:"cursorField"`
	Direction   int    `mapstructure:"direction"`
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
}

type PageInfo struct {
	StartCursor     interface{} `json:startCursor`
	EndCursor       interface{} `json:endCursor`
	HasNextPage     bool        `json:hasNextPage`
	HasPreviousPage bool        `json:hasPreviousPage`
}

type Connection struct {
	Nodes    []Map    `json:nodes` // TODO remove Map from func definition so other can replace with postgres, ...
	Edges    []Edge   `json:edges`
	PageInfo PageInfo `json:pageInfo`
}

type Edge struct {
	Node   Map         `json:node`
	Cursor interface{} `json:cursor`
}
