
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

pipeline = [
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
    nodes = data['nodes']
    
        
    for x in nodes:
        if ('messages' in x):
            x['_typename'] = 'MessageCampaign'
        elif ('posts' in x):
            x['_typename'] = 'PostCampaign'
        
    data['nodes'] = nodes
    return data

