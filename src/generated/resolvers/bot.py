
from tartiflette import Resolver
from .support import strip_nones, zip_pluck, select_keys
from operator import setitem
from funcy import omit

@Resolver('Query.bot')
async def resolve_query_bot(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    if not (where.get('user_id') == jwt.get('user_id')):
        raise Exception("guard `where.get('user_id') == jwt.get('user_id')` not satisfied")
    else:
        fields += ['ciao']
    
    collection = ctx['db']['bots']
    x = await collection.find_one(where)
    
    
    if fields:
        x = omit(x, fields)
    return x
