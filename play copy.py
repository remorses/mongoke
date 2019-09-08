
a = 'ciao'

guards = {
    a == 'ciao': ['ciao'],
    print('computed') or a == 'x' : [],
}


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


def filter_nodes_by_guard(nodes):
    excluded = []
    for x in nodes:
        if ok:
            excluded += fields
        else:
            continue
        yield omit(x, excluded)