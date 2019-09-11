from populate import populate_string

# scalars
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
""" for scalar in map(str, scalars)])
}}


type PageInfo {
    endCursor: AnyScalar
    startCursor: AnyScalar
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
    ${{query_name}}s(where: ${{type_name}}Where, cursorField: ${{type_name}}Fields, first: Int, last: Int, after: AnyScalar, before: AnyScalar): ${{type_name}}Connection
    ${{query_name}}(where: ${{type_name}}Where,): ${{type_name}}
}

type ${{type_name}}Connection {
    nodes: [${{type_name}}]
    pageInfo: PageInfo
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
   ${{relationName}}(where: ${{toType}}Where, cursorField: ${{toType}}Fields, first: Int, last: Int, after: AnyScalar, before: AnyScalar): ${{toType}}Connection
}
# if types don't have already the boilerplate i should write it now
'''

to_many_relation_boilerplate = '''
type ${{toType}}Connection {
    nodes: [${{toType}}]
    pageInfo: PageInfo
}

input ${{toType}}Where { 
    ${{'\\n    '.join([f'{field}: Where{scalar_name}' for field, scalar_name in fields])}}
}

enum ${{toType}}Fields {
    ${{'\\n    '.join([f'{field}' for field, _ in fields])}}
}
'''