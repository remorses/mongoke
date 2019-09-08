
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

def filter_nodes_by_guard(nodes):
    fields = []
    for x in nodes:
        try:
            
            yield omit(x, fields)
        except Exception:
            pass


pipeline: list = [
    {
        "$project": {
            "_id": 0,
            "username": 0
        }
    }
]

@Resolver('Query.campaigns')
async def resolve_query_campaigns(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    orderBy = args.get('orderBy', {'_id': 'ASC'}) # add default
    headers = ctx['request']['headers']
    jwt = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
    
    pagination = get_pagination(args)
    data = await connection_resolver(
        collection=ctx['db']['campaigns'], 
        where=where,
        orderBy=orderBy,
        pagination=pagination,
        pipeline=pipeline,
    )
    for x in data['nodes']:
        if ('messages' in x):
            x['_typename'] = 'MessageCampaign'
        elif ('posts' in x):
            x['_typename'] = 'PostCampaign'
        
    
    data['nodes'] = list(filter_nodes_by_guard(data['nodes']))
    return data

