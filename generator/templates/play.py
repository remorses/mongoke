
# where_filter, collection
@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = ${{repr_eval_dict(where_filter, '    ')}}
    x = await ctx['db']['${{collection}}'].find_one(where)
# ${{repr_guards_after_checks(guards_after, '    ')}}
# ${{repr_disambiguations(disambiguations, '    ')}}
    return x



# where_filter, collection
@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    relation_where = ${{repr_eval_dict(where_filter, '    ')}}
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)





@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    orderBy = args.get('orderBy', {'_id': 'ASC'}) # add default
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
${{repr_guards_before_checks(guards_before, '    ')}}
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

${{
"""
    nodes = []
    for x in data['nodes']:
""" if guards_after else """
    nodes = data['nodes']
"""
}}
${{filter_nodes_guards_after(guards_after, '        ')}}
${{
"""
    for x in nodes:
""" if disambiguations else ''
}}
${{repr_disambiguations(disambiguations, '        ')}}
    data['nodes'] = nodes
    return data
