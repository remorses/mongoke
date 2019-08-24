
import collections
from motor.motor_asyncio import AsyncIOMotorDatabase, AsyncIOMotorCollection
from tartiflette import Resolver
import pymongo
from pymongo import ASCENDING, DESCENDING
from typing import NamedTuple, Union
import typing

gt = '$gt'
lt = '$lt'
MAX_NODES = 20

parse_direction = lambda direction: ASCENDING if direction == 'ASC' else DESCENDING

async def connection_resolver(
        collection: AsyncIOMotorCollection, 
        where: dict,
        orderBy: dict, # needs to exist always at least one, the fisrst is the cursorField
        pagination: dict,
    ):
    first, last = pagination.get('first'), pagination.get('last'), 
    after, before = pagination.get('after'), pagination.get('before')
    first = min(MAX_NODES, first or 0)
    last = min(MAX_NODES, last or 0)

    sorting = [(field, parse_direction(direction)) for field, direction in orderBy.items()]
    cursorField = list(orderBy.keys())[0]

    if after and not (first or before):
        raise Exception('need `first` or `before` if using `after`')
    if before and not (last or after):
        raise Exception('need `last` or `after` if using `before`')
    if first and last:
        raise Exception('no sense using first and last together')

    if after != None and before != None:
        nodes = collection.find(
            {
                **where,
                cursorField: {
                    gt: after,
                    lt: before
                },
            },
            sort=sorting
        )
    elif after != None:
        nodes = collection.find(
            {
                **where,
                cursorField: {
                    gt: after,
                },
            },
            sort=sorting,
        )
    elif before != None:
        nodes = collection.find(
            {   
                **where,
                cursorField: {
                    lt: before
                },
            },
            sort=sorting,
        )
    else:
        nodes = collection.find(where, sort=sorting, )

    if first:
        nodes = nodes.limit(first + 1)
    elif last:
        toSkip = await collection.count_documents(where) - (last + 1)
        nodes = nodes.skip(max(toSkip, 0))

    nodes = await nodes.to_list(MAX_NODES)
    hasNext = None
    hasPrevious = None

    if first:
        hasNext = len(nodes) == (first + 1)
        nodes = nodes[:-1] if hasNext else nodes
        
    if last:
        hasPrevious = len(nodes) == last + 1
        nodes = nodes[1:] if hasPrevious else nodes

    end_cursor = nodes[-1][cursorField] if nodes else None
    start_cursor = nodes[0][cursorField] if nodes else None  
    return {
        'nodes': nodes,
        'pageInfo': {
            'endCursor': end_cursor,
            'startCursor': start_cursor,
            'hasNextPage': hasNext,
            'hasPreviousPage': hasPrevious,
        }
    }

def strip_nones(x: dict):
    result = {}
    for k, v in x.items():
        if not v == None:
            result[k] = v
    return result

