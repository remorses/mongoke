 
from tartiflette import Resolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem

pipeline: list = []

@Resolver('User.father')
async def resolve_user_father(parent, args, ctx, info):
    where = {
        "_id": {
            "$in":  parent['father_id'] 
        }
    }
    
    collection = ctx['db']['humans']
    x = await mongodb_streams.find_one(collection, match=where, pipeline=pipeline)
    
    # {{repr_disambiguations(disambiguations, '    ')
    return x
