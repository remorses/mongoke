package testutil

import (
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
