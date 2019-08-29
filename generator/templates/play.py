
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























@Resolver('Query.humans')
async def resolve_query_humans(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    orderBy = args.get('orderBy', {'_id': 'ASC'}) # add default
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
    if not (headers['user-id'] == where['_id'] or jwt_payload['user_id'] == 'ciao'):
        raise Exception("guard `headers['user-id'] == where['_id'] or jwt_payload['user_id'] == 'ciao'` not satisfied")
    else:
        fields += ['name', 'surname']
    
    pagination = get_pagination(args)
    data = await connection_resolver(
        collection=ctx['db']['humans'], 
        where=where,
        orderBy=orderBy,
        pagination=pagination,
    )
    def guard(x):
        return all([
            ,
        ])

    def disambiguate(x):
        if ('surname' in x):
            x['_typename'] = 'User'
        elif (x['type'] == 'guest'):
            x['_typename'] = 'Guest'
        return x
    data['nodes'] = [select_keys() for x in data['nodes'] if guard(x)]
    data['nodes'] = map(disambiguate, data['nodes'])
    data['nodes'] = list(data['nodes'])
    return data