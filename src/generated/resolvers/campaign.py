
from tartiflette import Resolver
from .support import strip_nones, zip_pluck, find_one, find
from operator import setitem
from funcy import omit

pipeline: list = [
    {
        "$project": {
            "_id": 0,
            "username": 0
        }
    }
]

@Resolver('Query.campaign')
async def resolve_query_campaign(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    
    collection = ctx['db']['campaigns']
    x = await find_one(collection, where, pipeline=pipeline)
    
    if ('messages' in x):
        x['_typename'] = 'MessageCampaign'
    elif ('posts' in x):
        x['_typename'] = 'PostCampaign'
    
    if fields:
        x = omit(x or dict(), fields)
    return x
