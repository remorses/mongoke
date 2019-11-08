# Relations
relations are defined in the configuration with the following shape
```
schema:
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
schema: |
    type Owner {
        _id: ObjectId
        email: String
    }
    type Pet {
        name: String
        owner_id: ObjectId
    }

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
schema: |
    type Owner {
        _id: ObjectId
        email: String
    }
    type Pet {
        name: String
        owner_id: ObjectId
        zoo_id: ObjectId
    }
    type Zoo {
        _id: ObjectId
        address: String
    }

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

