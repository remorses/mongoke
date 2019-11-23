
from tartiflette import Resolver, TypeResolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem
from funcy import omit

@TypeResolver('Human')
def resolve_type(result, context, info, abstract_type):
    x = result
    if (x['type'] == 'user'):
        return 'User'
    elif (x['type'] == 'guest'):
        return 'Guest'
    

pipeline: list = []

@Resolver('Query.human')
async def resolve_query_human(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].state.jwt_payload
    fields = []
    
    collection = ctx['db']['humans']
    x = await mongodb_streams.find_one(collection, match=where, pipeline=pipeline)
    if not (jwt['role'] == 'admin'):
        raise Exception("guard `jwt['role'] == 'admin'` not satisfied")
    else:
        fields += []
    if not (jwt['role'] == 'semi'):
        raise Exception("guard `jwt['role'] == 'semi'` not satisfied")
    else:
        fields += ['passwords', 'campaign_data']
    
    # {{repr_disambiguations(disambiguations, '    ')
    if fields:
        x = omit(x or dict(), fields)
    return x
