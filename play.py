from graphql import parse, Node, FieldDefinitionNode


s = '''
type Query {
    x: Int
    y: String
}
'''
definition: Node = parse(s).definitions[0]
print(dir(definition)) # DefinitionNode: 'description', 'directives', 'fields', 'interfaces', 'keys', 'kind', 'loc', 'name'
print(definition.name.value) # NameNode: 'value'
print(dir(definition.fields[0])) # FieldNode: 'arguments', 'description', 'directives', 'keys', 'kind', 'loc', 'name', 'type'
print(definition.fields[0].name.value)
print(dir(definition.fields[0].type)) # NamedTypeNode: 'keys', 'kind', 'loc', 'name'
print(definition.fields[0].type.name.value)