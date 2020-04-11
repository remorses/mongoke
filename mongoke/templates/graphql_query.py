from populate import populate_string

# searchable
general_graphql = '''

type Query {
    mongoke_version: String
}

enum Direction {
    ASC
    DESC
}

${{
''.join([f"""
input Where{scalar} {'{'}
    in: [{scalar}]
    nin: [{scalar}]
    eq: {scalar}
    neq: {scalar}
{'}'}
""" for scalar in map(str, sorted(searchables))])
}}


type PageInfo {
    startCursor: AnyScalar
    endCursor: AnyScalar
    hasNextPage: Boolean
    hasPreviousPage: Boolean
}

scalar Json
scalar ObjectId
scalar AnyScalar
'''

# fields
# query_name
graphql_query = '''
extend type Query {
    ${{query_name}}(
        where: ${{type_name}}Where,
    ): ${{type_name}}

    ${{query_name}}Nodes(
        where: ${{type_name}}Where, 
        cursorField: ${{type_name}}Fields, 
        first: Int, 
        last: Int, 
        after: AnyScalar, 
        before: AnyScalar,
        direction: Direction
    ): ${{type_name}}Connection!
}

type ${{type_name}}Connection {
    nodes: [${{type_name}}]!
    edges: [${{type_name}}Edge]!
    pageInfo: PageInfo!
}

type ${{type_name}}Edge {
    node: ${{type_name}}
    cursor: AnyScalar
}

input ${{type_name}}Where { 
    and: [${{type_name}}Where]
    or: [${{type_name}}Where]
    # $not: [${{type_name}}Where]
    ${{'\\n    '.join([f'{field}: Where{scalar_name}' for field, scalar_name in fields])}}
}

enum ${{type_name}}Fields {
    ${{'\\n    '.join([f'{field}' for field, _ in fields])}}
}
'''

######

to_one_relation = '''
extend type ${{fromType}} {
    ${{relationName}}: ${{toType}}
}
'''

to_many_relation = '''
extend type ${{fromType}} {
   ${{relationName}}(
       where: ${{toType}}Where, 
       cursorField: ${{toType}}Fields, 
       first: Int, 
       last: Int, 
       after: AnyScalar, 
       before: AnyScalar,
       direction: Direction
    ): ${{toType}}Connection!
}
'''

to_many_relation_boilerplate = '''
type ${{toType}}Connection {
    nodes: [${{toType}}]!
    edges: [${{toType}}Edge]!
    pageInfo: PageInfo!
}

type ${{toType}}Edge {
    node: ${{toType}}
    cursor: AnyScalar
}

input ${{toType}}Where { 
    ${{'\\n    '.join([f'{field}: Where{scalar_name}' for field, scalar_name in fields])}}
}

enum ${{toType}}Fields {
    ${{'\\n    '.join([f'{field}' for field, _ in fields])}}
}
'''