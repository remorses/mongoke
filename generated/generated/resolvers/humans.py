
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem
from funcy import omit

def filter_nodes_by_guard(nodes, fields):
    for x in nodes:
        try:
            if not (session['role'] == 'admin'):
                raise Exception("guard `session['role'] == 'admin'` not satisfied")
            else:
                fields += []
            
            yield omit(x or dict(), fields)
        except Exception:
            pass


map_fields_to_types = {
        "type": "String",
        "_id": "ObjectId",
        "name": "String",
        "surname": "String"
    }

pipeline: list = []

@Resolver('Query.humans')
async def resolve_query_humans(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    cursorField = args.get('cursorField', '_id')
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    if not (session['role'] == 'semi'):
        raise Exception("guard `session['role'] == 'semi'` not satisfied")
    else:
        fields += []
    
    pagination = get_pagination(args,)
    data = await connection_resolver(
        collection=ctx['db']['humans'], 
        where=where,
        cursorField=cursorField,
        pagination=pagination,
        scalar_name=map_fields_to_types[cursorField],
        pipeline=pipeline,
    )
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields))
    for x in data['nodes']:
        if (x['type'] == 'user'):
            x['_typename'] = 'User'
        elif (x['type'] == 'guest'):
            x['_typename'] = 'Guest'
        
    
    return data

