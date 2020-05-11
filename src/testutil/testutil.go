package testutil

import (
	"encoding/json"
	"testing"

	"github.com/graphql-go/graphql"
)

const MONGODB_URI = "mongodb://localhost/testdb"

func QuerySchema(t *testing.T, schema graphql.Schema, query string) interface{} {
	res := graphql.Do(graphql.Params{Schema: schema, RequestString: query})
	if res.Errors != nil && len(res.Errors) > 0 {
		t.Error(res.Errors[0])
		return nil
	}
	return res.Data
}

func ConvertToPlainMap(in interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(in)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func QuerySchemaShouldFail(t *testing.T, schema graphql.Schema, query string) error {
	res := graphql.Do(graphql.Params{Schema: schema, RequestString: query})
	if res.Errors != nil && len(res.Errors) > 0 {
		return res.Errors[0]
	}
	t.Fatal("query should have failed")
	return nil
}

var UserSchema = `
type User {
	name: String
	age: Int
}
`

var UserQuery1 = `
{
	findOneUser(where: {name: {eq: "dsf"}}) {
		name
		age
	}
}`

var UserQuery2 = `
{
	findManyUser(last: 1, where: {name: {eq: "sdfsdf"}}) {
	  nodes {
		name
		age
	  }
	  pageInfo {
		endCursor
		hasNextPage
		hasPreviousPage
		startCursor
	  }
	}
}`

var YamlConfig = `
jwt:
    algorithms: [H256]
    required: false

schema: |
    scalar Address
    scalar Url
    # scalar ObjectId
    type Task {
        _id: ObjectId
        address: Address
    }
    type WindowedEvent {
        value: Int
        timestamp: Int
    }
    type Guest {
        type: String
        _id: ObjectId
        name: String
    }
    enum Letter {
        a
        b
        c
    }
    type User {
        type: String
        _id: ObjectId
        name: String
        surname: String
        friends_ids: [ObjectId]
        url: Url
        letter: Letter
    }
    union Human = User | Guest

types:
    Task:
        collection: tasks
        exposed: false
    User:
        collection: users
        
    Human:
        collection: humans
        guards:
            -   expression: jwt['role'] == 'admin'
                when: after
            -   expression: jwt['role'] == 'semi'
                excluded: [passwords, campaign_data]
        disambiguations:
            User: x['type'] == 'user'
            Guest: x['type'] == 'guest'
    WindowedEvent:
        exposed: False
        collection: events
        pipeline:
            -   $group:
                    _id:
                        $substartct:
                            - $timestamp
                            - $mod: [$timestamp, 60000] # minute 60 * 1000
                    value:
                        $sum: $likes
            -   $project:
                    _id: 0
                    value: 1
                    timestamp: $_id

relations:
    -   from: Task
        to: WindowedEvent
        relation_type: to_many
        field: events
        where: {}
    -   from: User
        to: User
        relation_type: to_many
        field: friends
        where: {}
    -   from: User
        to: WindowedEvent
        relation_type: to_many
        field: likes_over_time
        where:
            bot_id:
                in: ["${{ parent['_id'] }}"]
            type:
                eq: like
    -   from: User
        to: Human
        field: father
        relation_type: to_one
        where:
            _id:
                in: ["${{ parent['father_id'] }}"]


`

const IntrospectionQuery = `
  query IntrospectionQuery {
    __schema {
      queryType { name }
      mutationType { name }
      subscriptionType { name }
      types {
        ...FullType
      }
      directives {
        name
        description
        locations
        args {
          ...InputValue
        }
        # deprecated, but included for coverage till removed
        onOperation
        onFragment
        onField
      }
    }
  }
  fragment FullType on __Type {
    kind
    name
    description
    fields(includeDeprecated: true) {
      name
      description
      args {
        ...InputValue
      }
      type {
        ...TypeRef
      }
      isDeprecated
      deprecationReason
    }
    inputFields {
      ...InputValue
    }
    interfaces {
      ...TypeRef
    }
    enumValues(includeDeprecated: true) {
      name
      description
      isDeprecated
      deprecationReason
    }
    possibleTypes {
      ...TypeRef
    }
  }
  fragment InputValue on __InputValue {
    name
    description
    type { ...TypeRef }
    defaultValue
  }
  fragment TypeRef on __Type {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
                ofType {
                  kind
                  name
                }
              }
            }
          }
        }
      }
    }
  }
`
