
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

@Resolver('Query.humans')
async def resolve_query_humans(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    orderBy = args.get('orderBy', {'_id': 'ASC'}) # add default
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
    if not (session['role'] == 'semi'):
        raise Exception("guard `session['role'] == 'semi'` not satisfied")
    else:
        fields += []
    
    pagination = get_pagination(args)
    data = await connection_resolver(
        collection=ctx['db']['users'], 
        where=where,
        orderBy=orderBy,
        pagination=pagination,
    )

    nodes = []
    for x in data['nodes']:

        if not (session['role'] == 'admin'):
            pass
        else:
            own_fields = fields + []
            if own_fields:
                x = select_keys(lambda k: k in fields, x)
            nodes.append(x)
        

    for x in nodes:

        if (x['type'] == 'user'):
            x['_typename'] = 'User'
        
        elif (x['type'] == 'guest'):
            x['_typename'] = 'Guest'
        
    data['nodes'] = nodes
    return data
