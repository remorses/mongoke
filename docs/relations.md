# Relations
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

