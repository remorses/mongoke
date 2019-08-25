
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

@Resolver('Human.friends')
async def resolve_human_friends(parent, args, ctx, info):
    relation_where = {
        "_id": {
            "$in": parent['friends_ids']
        }
    }
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)
    