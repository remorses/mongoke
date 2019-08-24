
from tartiflette import Resolver
from .support import strip_nones, connection_resolver
from operator import setitem

@Resolver('Query.human')
async def resolve_query_human(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['request']['headers']
    jwt_payload = ctx['req'].jwt_payload # TODO i need to decode jwt_payload and set it in req in a middleware
    if not (True):
        raise Exception('guard True not satisfied')
    collection=ctx['db']['collection}']
    x = collection.find_one(where)
    
    if not True:
        raise Exception('guard True not satisfied')
    return x
