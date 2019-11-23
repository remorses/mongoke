
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem
from funcy import omit

def filter_nodes_by_guard(nodes, fields, jwt):
    for x in nodes:
        try:
            if not (jwt['role'] == 'admin'):
                raise Exception("guard `jwt['role'] == 'admin'` not satisfied")
            else:
                fields += []
            if not (jwt['role'] == 'semi'):
                raise Exception("guard `jwt['role'] == 'semi'` not satisfied")
            else:
                fields += ['passwords', 'campaign_data']
            
            yield omit(x or dict(), fields)
        except Exception:
            pass


map_fields_to_types = {
        "type": "String",
        "_id": "ID",
        "name": "String",
        "surname": "String"
    }

pipeline: list = []

@Resolver('Query.humans')
async def resolve_query_humans(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    cursorField = args.get('cursorField',) or ('_id' if '_id' in map_fields_to_types else list(map_fields_to_types.keys())[0])
    headers = ctx['req'].headers
    jwt = ctx['req'].state.jwt_payload
    fields = []
    
    pagination = get_pagination(args,)
    data = await connection_resolver(
        collection=ctx['db']['humans'], 
        where=where,
        cursorField=cursorField,
        pagination=pagination,
        scalar_name=map_fields_to_types[cursorField],
        pipeline=pipeline,
    )
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields, jwt=jwt))
    # {{repr_many_disambiguations(disambiguations, '    ') if disambiguations else ''
    return data

