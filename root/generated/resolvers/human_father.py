 
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys
from operator import setitem

@Resolver('Human.father')
async def resolve_human_father(parent, args, ctx, info):
    where = {
        "_id": {
            "$in": parent['father_id']
        }
    }

    x = await ctx['db']['users'].find_one(where)


    return x