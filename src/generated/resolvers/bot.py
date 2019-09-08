
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
    
    collection = ctx['db']['bots']
    x = await collection.find_one(where)
    
    
    if fields:
        x = omit(x, fields)
    return x
