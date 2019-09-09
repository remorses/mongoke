
from tartiflette import Resolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem
from funcy import omit

pipeline: list = [
    {
        "$set": {
            "username": "fucku"
        }
    }
]

@Resolver('Query.bot')
async def resolve_query_bot(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    
    collection = ctx['db']['bots']
    x = await mongodb_streams.find_one(collection, where, pipeline=pipeline)
    
    
    if fields:
        x = omit(x or dict(), fields)
    return x
