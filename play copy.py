
a = 'ciao'

guards = {
    a == 'ciao': ['ciao'],
    print('computed') or a == 'x' : [],
}

excluded = []
if True in guards:
    excluded += guards[True]
else:
    raise Exception('no expression satisfied')

node = omit(node, excluded)

disambiguations = {
    'ciao' in a: 'User'
}

node = {}
node['_typename'] = disambiguations[True]