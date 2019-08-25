from populate import populate_string


general_graphql = '''
enum Direction {
    ASC
    DESC
}

input WhereInput {
    $in: [NumberOrString]
    $nin: [NumberOrString]
    $eq: NumberOrString
    $neq: NumberOrString
    $gte: NumberOrString
    $lte: NumberOrString
    $gt: NumberOrString
    $lt: NumberOrString
}

type PageInfo {
    endCursor: Int
    startCursor: Int
    hasNextPage: Boolean
    hasPreviousPage: Boolean
}

scalar Json
scalar Url
scalar ObjectId
scalar NumberOrString

'''

# fields
# query_name
graphql_query = '''
extend type Query {
    ${{query_name}}s(where: ${{type_name}}Where, orderBy: ${{type_name}}OrderBy,): ${{type_name}}Connection
    ${{query_name}}(where: ${{type_name}}Where,): ${{type_name}}
}

type ${{type_name}}Connection {
    nodes: [${{type_name}}]
    pageInfo: PageInfo
}

input ${{type_name}}Where { 
    ${{'\\n    '.join([f'{field}: WhereInput' for field in fields])}}
}

input ${{type_name}}OrderBy {
    ${{'\\n    '.join([f'{field}: Direction' for field in fields])}}
}
'''

######

to_one_relation = '''
extend type ${{fromType} {
    ${{relationName}}: ${{toType}}
}
'''

to_many_relation = '''
extend type ${{fromType} {
   ${{relationName}}(where: ${{toType}}Where, orderBy: ${{toType}}OrderBy): ${{toType}}Connection
}
# if types don't have already the boilerplate i should write it now
'''

to_many_relation_boilerplate = '''
type ${{toType}}Connection {
    nodes: [${{type_name}}]
    pageInfo: PageInfo
}

input ${{toType}}Where { 
    ${{'\\n    '.join([f'{field}: WhereInput' for field in fields])}}
}

input ${{toType}}OrderBy {
    ${{'\\n    '.join([f'{field}: Direction' for field in fields])}}
}
'''