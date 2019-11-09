
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem
from funcy import omit

def filter_nodes_by_guard(nodes, fields, jwt):
    for x in nodes:
        try:
            
            yield omit(x or dict(), fields)
        except Exception:
            pass


map_fields_to_types = {
        "value": "Int",
        "timestamp": "Int"
    }

pipeline: list = [
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

@Resolver('Task.events')
async def resolve_task_events(parent, args, ctx, info):
    relation_where = {}
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)
    cursorField = args.get('cursorField',) or ('_id' if '_id' in map_fields_to_types else list(map_fields_to_types.keys())[0])
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    
    pagination = get_pagination(args,)
    data = await connection_resolver(
        collection=ctx['db']['events'], 
        where=where,
        cursorField=cursorField,
        pagination=pagination,
        scalar_name=map_fields_to_types[cursorField],
        pipeline=pipeline,
    )
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields))
    # {{repr_many_disambiguations(disambiguations, '    ') if disambiguations else ''
    return data
