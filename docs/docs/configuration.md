---
route: /docs/configuration
name: Configuration
---

## Configuration

Mongoke defines its entire configuration in a yaml file that can be used to generate the entire graphql server.
This configuration can be used inside the docker image in the default path `/config.yml`.

```yaml
Configuration:
    schema?: Str
    schema_url?: Url
    schema_path?: Str
    types:
        ...:
            collection: Str
            exposed?: Bool
            pipeline?: [Any]
            disambiguations?:
                ...: Str
            guards?: [
                expression: Str
                excluded?: [Str]
                when?: "after" | "before"
            ]
    relations?: [
        from: Str
        to: Str
        relation_type: "to_many" | "to_one"
        field: Str
        where: Any
    ]
    jwt?:
        secret?: Str
        header_name?: Str # default is "Authorization"
        header_scheme?: Str # default is "Bearer"
        required?: Bool
        algorithms?: ["H256" | "HS512" | "HS384" | "RS256" | "RS384" | "RS512" | "ES256" | "ES384" | "ES521" | "ES512" | "PS256" | "PS384" | "PS512"]

Url: Str
```


## Types

Individual types config defined as an object where keys are the type names and values are the type configuration.

```
types:
    ...:
        collection: Str
        exposed?: Bool
        pipeline?: [Any]
        guards?: [
            expression: Str
            when?: "after" | "before"
            excluded?: [Str]
        ]
        disambiguations?:
            ...: Str
```

### collection

Defines to what collection the type is associated with

### exposed

Defines if the type is exposed to graphql, useful when you whant to use certain types only as relations

### pipeline

custom mongodb pipeline to execute during the database query

### disabiguations

necessary when querying a union type, to determine the actual type.
it is an object where keys are type names and values are expressions.
The expressions are evaluated until one is found true and the right \_\_type is applied

## Guards

List of expressions to limit the access of the type fields to only certain users, based on jwt payload and the document data.

```
guards?: [
    expression: Str
    excluded?: [Str]
    when?: "after" | "before"
]
```

### when

decides if you want to evaluate the expression before or after querying the database, if you use before you save resources but have access to only the user jwt (if any) and not to the document to decide if user is authorized

### expression

python expression that can evaluate to true if you want to give user access to the type, expression is evaluated in python and has access to

-   x: the current document, available only if using when=after
-   jwt: the user jwt payload, can contain whatever you put inside it, by default extracted from the Authorization header and not verified.

### excluded

By default the guards give access to all the document fields, you can limit the fileds you give access to by putting them inside `exclude`.
To implement different levels of authorization with access to different fields you can use many guards where the most protected is the first so that the evaluation stops at the weakest permissions required possible.

## Relations

Defined as a list of configurations to add connections between types.

```
relations?: [
    from: Str
    to: Str
    relation_type: "to_many" | "to_one"
    field: Str
    where: Any # the mongodb query
]
```

### from

The type where the relation's field is added

### to

The type the relation leads to

### field

The field added to the `from` type to connect the `to` type

### relation_type

if "to_one" the field in graphql will be a simple type reference and can be queried with

```gql
{
    owner {
        email
        pet {
            name
        }
    }
}
```

If "to_many" the field will resolve to a connection and can be queried like this

```gql
{
    zoo {
        pets(first: 10) {
            nodes {
                name
            }
        }
    }
}
```

### where

The mongodb where query to find the related documents, you can evaluate custom python code inside the \${{ }} and have access to parent: the `from` document as a python dict.
The code inside \${{ }} will be evaluated during every query that needs the relation and the evaluation result will be used to query the `to` collection.

## Jwt configuration

Configure how to handle jwt authentication, by default the jwt is not verified, to verify it add the `secret` field with the secret used to sign the jwt. You can require a jwt for all the query fields adding the `required` field.

```
    jwt?:
        secret?: Str
        header_name?: Str # default is "Authorization"
        header_scheme?: Str # default is "Bearer"
        required?: Bool
        algorithms?: ["H256" | "HS512" | "HS384" | "RS256" | "RS384" | "RS512" | "ES256" | "ES384" | "ES521" | "ES512" | "PS256" | "PS384" | "PS512"]
```

### required

if specified, only users with jwt signed with the right secret have access to the resources, needs secret to work.
By default the secret is not required and not verified.

### secret

Used when required is present to check if jwt is signed

### algorithms

A list of algotihtm to decode the jwt, to see the full list chech the python pyJwt library
