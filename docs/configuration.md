## Configuration

Mongoke defines its entire configuration in a yaml file that can be used to generate the entire graphql server.
This configuration can be used inside the docker image in the default path `/config.yml`, read more about how to use mongoke with docker here.
The configuration has the following shape:
```yml
Configuration:
    db_url?: /mongodb://.*/
    skema?: Str
    skema_url?: Url
    skema_path?: Str
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
        secret: Str
        algorithms: ["H256"]

Url: Str
```

## types

```
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
```

