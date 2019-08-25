
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys
from operator import setitem
from funcy import select_keys

@Resolver('Query.human')
async def resolve_query_human(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload and set it in req in a middleware
    fields = []
    if not (headers['user-id'] == where['_id'] or jwt_payload['user_id'] == 'ciao'):
        raise Exception("guard `headers['user-id'] == where['_id'] or jwt_payload['user_id'] == 'ciao'` not satisfied")
    else:
        fields += ['name', 'surname']
    
    collection = ctx['db']['humans']
    x = collection.find_one(where)

    if ('surname' in x):
        x['_typename'] = 'User'
    
    elif (x['type'] == 'guest'):
        x['_typename'] = 'Guest'
    
    if fields:
        x = select_keys(lambda k: k in fields, x)
    return x