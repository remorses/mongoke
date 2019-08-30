 
from tartiflette import Resolver
from .support import strip_nones, zip_pluck, select_keys
from operator import setitem

@Resolver('Bot.user')
async def resolve_bot_user(parent, args, ctx, info):
    where = {
        "_id":  parent['_id'] 
    }

    x = await ctx['db']['campaigns'].find_one(where)


    return x