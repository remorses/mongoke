 
from tartiflette import Resolver
from .support import strip_nones, zip_pluck, find_one, find
from operator import setitem

pipeline: list = []

@Resolver('Bot.user')
async def resolve_bot_user(parent, args, ctx, info):
    where = {
        "_id":  parent['_id'] 
    }
    
    collection = ctx['db']['campaigns']
    x = await find_one(collection, where, pipeline=pipeline)
    
    
    return x
