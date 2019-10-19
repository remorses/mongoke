## Database types

mongodb is a schemaless database, you can put whatever you like inside it, but to generate automatically an api the configuration must specify a schema that describes the expected documents shape for every exposed collection.
To do this in mongoke we use a language made just for types, skema.

## Skema
[Skema](https://github.com/remorses/skema) is a DSL born from the necessity of defining domain types in one place for every language used in a multi language microservice architecture. One of the languages skema compiles to is graphql.
Here is an example bit of skema and his generated gralphql types:
```yml
BlogPost:
    title: Str
    description: Str
    user:
        username: Str
        age: Int
    tags: [
        name: Str
        id: Int
    ]
```
```gql
type BlogPost {
    title: String
    description: String
    user: User
    tags: [Tags]
}

type Tags {
    name: String
    id: Int
}

type User {
    username: String
    age: Int
}
```
As you can see skema supports nested object types and an easier list definition.
You can read more about skema [here](https://github.com/remorses/skema).

## Built in types
Skema has some built in types:
- Str
- Bool
- Int
- Float
- Any

In mongoke there are some additional scalars implemented by default:
- DateTime
- ObjectId
- Date
- Time

