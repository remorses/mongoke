
from tartiflette import Resolver, TypeResolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem
from funcy import omit



pipeline: list = []

@Resolver('Query.user')
async def resolve_query_user(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    
    collection = ctx['db']['users']
    x = await mongodb_streams.find_one(collection, match=where, pipeline=pipeline)
    
    # {{repr_disambiguations(disambiguations, '    ')
    if fields:
        x = omit(x or dict(), fields)
    return x
