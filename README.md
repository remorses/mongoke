Todo:
- unit tests for the connection_resolver
- integration tests for all the resolver types
- integration tests for the relations
- cursor must be obfuscated in connection, (also after and before are string so it is a must)
- ~~add pipelines feature to all resolvers (adding a custom find and find_one made with aggregate)~~
- ~~add the $ to the where input fields inside resolvers (in must be $in, ...)~~
- ~~remove strip_nones after asserting v1 works~~

Low priority
- add verify the jwt with the secret if provided
- ~~add schema validation to the configuration~~
- add subscriptions
- add edges to make connection type be relay compliant 
- better performance of connection_resolver removing the $skip and $count
- add a dataloader for single connections


connection_resolver coercion

add cursorField argument,
in every resolver create a map from cursorField to scalar typename
map = {
    _id: ObjectId,
    name: String
}

then at the time of get_pageination pass as argument the scalar type, doing

def get_cursor_coercer(info):
    field = args.get('cursorField', '_id')
    scalar_name = map_fields_to_types[field]
    scalars = info.schema._scalar_definitions
    return scalars[scalar_name].input_coercer

get_pagination(args, get_cursor_coercer(info))

def get_pagination(args, coerce):
    after = args.get('after')
    before = args.get('before')
    return {
        'after': after and coerce(after),
        'before': before and coerce(before),
        'first': args.get('first'),
        'last': args.get('last'),
    }
