## Relations
relations are defined in the configuration with the following shape
```
skema:
    ...
types: 
    ...
relations?: [
    from: Str
    to: Str
    relation_type: "to_many" | "to_one"
    field: Str
    where: Any
]
```
-------

One example of a configuration with a one `to_one` relation:
```
skema: |
    Owner:
        _id: ObjectId
        email: Str
    Pet:
        name: Str
        owner_id: ObjectId

types:
    Owner:
        collection: owners
    Pet:
        collection: pets

relations:
    -   field: pet
        relation_type: to_one
        from: Owner
        to: Pet
        where: { "owner_id": ${{ parent['_id'] }}}
```
An example of a `to_many` relation:
```
skema: |
    Owner:
        _id: ObjectId
        email: Str
    Pet:
        name: Str
        owner_id: ObjectId
        zoo_id: ObjectId
    Zoo:
        _id: ObjectId
        address: Str

types:
    Owner:
        collection: owners
    Pet:
        collection: pets
    Zoo:
        collection: zoos


relations:
    -   field: pets
        relation_type: to_many
        from: Zoo
        to: Pet
        where: { "zoo_id": ${{ parent['_id'] }}}
```


** from **
The type where the relation's field is added
** to **
The type the relation leads to
**field**
The field added to the `from` type to connect the `to` type
**relation_type**
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
**where**
The mongodb where query to find the related documents, you can evaluate custom python code inside the ${{ }} and have access to parent: the `from` document as a python dict.
The code inside ${{ }} will be evaluated during every query that needs the relation and the evaluation result will be used to query the `to` collection.


