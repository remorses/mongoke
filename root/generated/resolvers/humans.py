
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys
from operator import setitem

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
    
    pagination = {
        'after': args.get('after'),
        'before': args.get('before'),
        'first': args.get('first'),
        'last': args.get('last'),
    }
    data = await connection_resolver(
        collection=ctx['db']['humans'], 
        where=where,
        orderBy=orderBy,
        pagination=pagination,
    )
    # User: 'surname' in x
    # Guest: x['type'] == 'guest'


    nodes = data['nodes']



    for x in nodes:

        if ('surname' in x):
            x['_typename'] = 'User'
        
        elif (x['type'] == 'guest'):
            x['_typename'] = 'Guest'
        
    data['nodes'] = nodes
    return data
