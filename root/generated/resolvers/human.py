
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
    if not (session['role'] == 'semi'):
        raise Exception("guard `session['role'] == 'semi'` not satisfied")
    else:
        fields += []
    
    collection = ctx['db']['users']
    x = collection.find_one(where)
    if not (session['role'] == 'admin'):
        raise Exception("guard `session['role'] == 'admin'` not satisfied")
    else:
        fields += []
    
    if (x['type'] == 'user'):
        x['_typename'] = 'User'
    
    elif (x['type'] == 'guest'):
        x['_typename'] = 'Guest'
    
    if fields:
        x = select_keys(lambda k: k in fields, x)
    return x