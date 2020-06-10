
getUserRelayConnection

findOneUser
findManyUser

insertOneUser
insertManyUser

updateOneUser
updateManyUser

deleteOneUser
deleteManyUser

upsertOneUser
upsertManyUser




functions

/mongoke
`NewMongokeFromConfig(config interface{})`
read a yaml configuration, creates the config and calls 
- `mongoke.generateTypeMap()`
- `mongoke.generateSchema()`


`Mongoke{config Config}`

/types
`mongoke.generateTypeMap()`
for every type in the model generates a where argument, connection, scalarFieldsEnum and stores them in the `mongoke.typeMap`
every type here must be a reference

/schema
`mongoke.generateSchema()`
for every type attach a field in the final schema

/fields
`mongoke.findOneField(object graphql.Object) graphql.Field`
`mongoke.findManyField(object)`
generates a field with resolver for the given type, takes other necessary types from the typeMap



/support

`makeMongodbMatch(map[string]interface{})`
replaces the `eq` `neq` fields with `$eq` `$neq`

`connectionResolver(collection, first, last, after, ...)`
get the requested segment from the database, i also have to redeclare the Connection struct with right json names


