
from tartiflette import Resolver
from .support import strip_nones, zip_pluck, select_keys
from operator import setitem
from funcy import select_keys

@Resolver('Query.bot')
async def resolve_query_bot(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['request']['headers']
    jwt = ctx['req'].jwt_payload # TODO i need to decode jwt_payload and set it in req in a middleware
    fields = []
    if not (jwt['_id'] == where['_id']):
        raise Exception("guard `jwt['_id'] == where['_id']` not satisfied")
    else:
        fields += []
    
    collection = ctx['db']['bots']
    x = collection.find_one(where)
    
    
    if fields:
        x = select_keys(lambda k: k in fields, x)
    return x
