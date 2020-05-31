package fakedata

import (
	"strings"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/remorses/mongoke/src/testutil"
)

func TestGenerate(t *testing.T) {
	tests := map[string]struct {
		typeDefs       string
		typeName       string
		expectedFields []string
		expectedError  string
	}{
		"simple": {
			typeName: "User",
			typeDefs: `
				type User {
					name: String
					age: Int
				}
			`,
			expectedFields: []string{
				"name",
				"age",
			},
		},
		"nested": {
			typeName: "User",
			typeDefs: `
				type User {
					name: String
					age: Int
					address: Address
				}
				type Address {
					city: String
					state: String
				}
			`,
			expectedFields: []string{
				"name",
				"age",
				"address.city",
				"address.state",
			},
		},
		"nested with enum": {
			typeName: "User",
			typeDefs: `
				type User {
					name: String
					age: Int
					address: Address
				}
				type Address {
					state: State
				}
				enum State {
					CA
					LA
					NY
				}
			`,
			expectedFields: []string{
				"name",
				"age",
				"address.state",
			},
		},
		"using custom scalar": {
			typeName: "User",
			typeDefs: `
				scalar ObjectId
				type User {
					_id: ObjectId
					name: String
					age: Int
				}
			`,
			expectedFields: []string{
				"name",
				"age",
				"_id",
			},
		},
		"union type": {
			typeName: "Account",
			typeDefs: `
				union Account = User | Guest
				type Guest {
					anonymous: Boolean
				}
				type User {
					name: String
					age: Int
				}
			`,
			expectedFields: []string{
				"name",
				"age",
				"anonymous",
			},
		},
		"interface type": {
			typeName: "User",
			typeDefs: `
				interface User {
					name: String
					age: Int
				}
			`,
			expectedFields: []string{
				"name",
				"age",
			},
		},
	}
	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			fakeData, err := NewFakeData(NewFakeDataParams{
				typeDefs: test.typeDefs,
			})
			if err != nil {
				t.Error(err)
			}
			x, err := fakeData.Generate(test.typeName)
			if err != nil {
				t.Error(err)
			}
			json := testutil.Pretty(x)
			t.Log(json)
			for _, field := range test.expectedFields {
				parts := strings.Split(field, ".")
				_, dt, _, err := jsonparser.Get([]byte(json), parts...)
				if err != nil {
					t.Error(err)
				}
				if dt == jsonparser.NotExist {
					t.Error(field + " not present")
				}
			}
		})
	}
}
