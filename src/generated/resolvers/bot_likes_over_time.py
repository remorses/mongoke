
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

pipeline = [
    {
        "$group": {
            "_id": {
                "$substartct": [
                    "$timestamp",
                    {
                        "$mod": [
                            "$timestamp",
                            60000
                        ]
                    }
                ]
            },
            "value": {
                "$sum": "$likes"
            }
        }
    },
    {
        "$project": {
            "_id": 0,
            "value": 1,
            "timestamp": "$_id"
        }
    }
]

@Resolver('Bot.likes_over_time')
async def resolve_bot_likes_over_time(parent, args, ctx, info):
    relation_where = {
        "bot_id": {
            "$in":  parent['_id'] 
        },
        "type": "like"
    }
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)
    
