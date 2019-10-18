# Generated queries

every type defined in the configurations and associated with a collection will generate a graphql query for a single document and for many documents in the form of a connection.

For example having the following type the following graphql types will be generated:
```
User:
    type: "user"
    _id: ID
    name: Str
    surname: Str
    friends_ids: [ID]
```
```graphql
extend type Query {
    user(
        where: UserWhere,
    ): User

    users(
        where: UserWhere, 
        cursorField: UserFields, 
        first: Int, 
        last: Int, 
        after: AnyScalar, 
        before: AnyScalar
    ): UserConnection

}

type UserConnection {
    nodes: [User]
    pageInfo: PageInfo
}

input UserWhere { 
    and: [UserWhere]
    or: [UserWhere]
    type: WhereString
    _id: WhereID
    name: WhereString
    surname: WhereString
}

enum UserFields {
    type
    _id
    name
    surname
}
type PageInfo {
    endCursor: AnyScalar
    startCursor: AnyScalar
    hasNextPage: Boolean
    hasPreviousPage: Boolean
}
```

Some examples pf queries for these types:
```gql

{
    user(where: {name: {eq: "jon"}) {
        _id
        name
        surname
    }
}
```

```gql

{
    users(first: 20, after: "Micky", cursosorField: name) {
        nodes {
            _id
            name
            surname
        }
        pageInfo {
            endCursor
        }
    }
}
```


The relations are similar, for a to_one relation:
```graphql

extend type User {
    father: User
}
```
For to_many relations:
```gql
extend type User {
   friends(
       where: UserWhere, 
       cursorField: UserFields, 
       first: Int, 
       last: Int, 
       after: AnyScalar, 
       before: AnyScalar
    ): UserConnection
}
```
