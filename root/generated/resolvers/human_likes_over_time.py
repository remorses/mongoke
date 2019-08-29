
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

@Resolver('Human.likes_over_time')
async def resolve_human_likes_over_time(parent, args, ctx, info):
    relation_where = {
        "bot_id": {
            "$in":  parent['_id'] 
        },
        "type": "like"
    }
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)
    