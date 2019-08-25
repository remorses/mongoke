from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck
from operator import setitem
from funcy import select_keys

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload and set it in req in a middleware
    fields = []
${{
'\\n'.join([f"""
    if not ({expr}):
        raise Exception('guard {expr} not satisfied')
    else:
        fields += {fields}
"""
for expr, fields in zip_pluck(guards_before, ['expression', 'fields'])])
}}
    collection=ctx['db']['${{collection}']
    x = collection.find_one(where)
${{
'\\n'.join([f"""
    if not ({expr}):
        raise Exception('guard {expr} not satisfied')
    else:
        fields += {fields or []}
"""
for expr, fields in zip_pluck(guards_after, ['expression', 'fields'])])
}}
    if fields:
        x = select_keys(lambda k: k in fields, x)
${{
'\\n'.join([
f"""
    if ({expr}):
        x['_typename'] = '{typename}'
""" 
for typename, expr in zip_pluck(disambiguations, ['type_name', 'expression'])])
}}
    return x




























from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck
from operator import setitem

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    orderBy = args.get('orderBy', {'_id': 'ASC'}) # add default
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
${{
'\\n'.join([f"""
    if not ({expr}):
        raise Exception('guard {expr} not satisfied')
    else:
        fields += {fields}
"""
for expr, fields in zip_pluck(guards_before, ['expression', 'fields'])])
}}
    pagination = {
        'after': args.get('after'),
        'before': args.get('before'),
        'first': args.get('first'),
        'last': args.get('last'),
    }
    data = await connection_resolver(
        collection=ctx['db']['${{collection}}'], 
        where=where,
        orderBy=orderBy,
        pagination=pagination,
    )
    # User: 'surname' in x
    # Guest: x['type'] == 'guest'
    nodes = []
    for x in data['nodes']:
${{
'\\n'.join([f"""
        if not ({expr}):
            pass
        else:
            own_fields = fields + {fields or []}
            if own_fields:
                x = select_keys(lambda k: k in fields, x)
            nodes.append(x)
"""
for expr, fields in zip_pluck(guards_after, ['expression', 'fields'])])
}}
    data['nodes'] = nodes
${{
'\\n'.join([
f"""
for x in data['nodes']:
    if ({expr}):
        x['_typename'] = '{typename}'
""" 
for typename, expr in zip_pluck(disambiguations, ['type_name', 'expression'])])
}}
    return data

nodes = []
for x in data['nodes']:
    if not ({expr}):
        pass
    else:
        own_fields = fields + {fields or []}
        if own_fields:
            x = select_keys(lambda k: k in fields, x)
        nodes.append(x)
data['nodes'] = nodes
for x in data['nodes']:
    if ({expr}):
        x['_typename'] = '{typename}'
