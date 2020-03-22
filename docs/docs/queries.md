---
route: /docs/queries
title: queries
---

Every type defined in the configurations and associated with a collection will generate a graphql query for a single document and for many documents.

The best way to explore the queries shape is follow the quickstart guide and open graphiql to explore the possible queries.

For example having the following type you can do the following queries

``` 
User:
    type: "user"
    _id: ID
    name: Str
    surname: Str
    friends_ids: [ID]
```

Some examples pf queries for this type:

``` gql

{
    user(where: {name: {eq: "jon"}) {
        _id
        name
        surname
    }
}
```

Every type generates a where argument where you can query the mongodb database with the `eq` , `in` , `nin` and other mongodb operators.

``` gql
{
    users(first: 20, after: "Micky", cursosorField: name) {
        nodes {
            _id
            name
            surname
        }
        pageInfo {
            endCursor
            startCursor
            hasNextPage
            hasPreviousPage
        }
    }
}
```

The connections have additional arguments to handle pagination, the documents are always sorted ascending on the \_id field if present, you can change the sorting field with the `cursorField` argument.
The pageInfo field returns the information to handle pagination, the endCursor and startCursor fields can be any scalar type based on the cursorField argument, they are not obfuscated to make it easier to see what is happening inside your app.

The generated graphql is below.

``` graphql
extend type Query {
    user(where: UserWhere): User

    users(
        where: UserWhere
        cursorField: UserFields
        first: Int
        last: Int
        after: AnyScalar
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

The relations are similar, for a to_one relation:

``` graphql
extend type User {
    father: User
}
```

For to_many relations:

``` gql
extend type User {
    friends(
        where: UserWhere
        cursorField: UserFields
        first: Int
        last: Int
        after: AnyScalar
        before: AnyScalar
    ): UserConnection
}
```

