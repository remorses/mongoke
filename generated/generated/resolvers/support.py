
import collections
from motor.motor_asyncio import AsyncIOMotorDatabase, AsyncIOMotorCollection
import mongodb_streams
from tartiflette import Resolver
import pymongo
from pymongo import ASCENDING, DESCENDING
from typing import NamedTuple, Union
import typing
from funcy import pluck, select_keys, omit
from ..scalars import scalar_classes


gt = '$gt'
lt = '$lt'
MAX_NODES = 20
DEFAULT_NODES_COUNT = 10

INPUT_COERCERS = {
    'String': str,
    'Int': int,
    'Float': float,
    'Bool': bool,
    **{scalar.name: scalar._implementation.coerce_input for scalar in scalar_classes},
}

OUTPUT_COERCERS = {
    'String': str,
    'Int': int,
    'Float': float,
    'Bool': bool,
    **{scalar.name: scalar._implementation.coerce_output for scalar in scalar_classes},
}

def zip_pluck(d, *keys):
    return zip(*[pluck(k, d) for k in keys])

def get_pagination(args,):
    after = args.get('after')
    before = args.get('before')
    return {
        'after': after,
        'before': before,
        'first': args.get('first'),
        'last': args.get('last'),
    }



async def connection_resolver(
    collection: AsyncIOMotorCollection,
    where: dict,
    cursorField,  # needs to exist always at least one, the fisrst is the cursorField
    pagination: dict,
    scalar_name,
    pipeline=[]
):
    first, last = pagination.get('first'), pagination.get('last'),
    after, before = pagination.get('after'), pagination.get('before')
    if after:
        after = INPUT_COERCERS[scalar_name](after)
    if before:
        before = INPUT_COERCERS[scalar_name](before)

    first = min(MAX_NODES, first or 0)
    last = min(MAX_NODES, last or 0)

    if not first and not last:
        if after:
            first = DEFAULT_NODES_COUNT
        elif before:
            last = DEFAULT_NODES_COUNT
        else:
            first = DEFAULT_NODES_COUNT

    if after and not (first or before):
        raise Exception('need `first` or `before` if using `after`')
    if before and not (last or after):
        raise Exception('need `last` or `after` if using `before`')
    if first and last:
        raise Exception('no sense using first and last together')

    args: dict = dict()

    if after != None and before != None:
        args.update(dict(
            match={
                **where,
                cursorField: {
                    gt: after,
                    lt: before
                },
            },
        ))
    elif after != None:
        args.update(dict(
            match={
                **where,
                cursorField: {
                    gt: after,
                },
            },
        ))
    elif before != None:
        args.update(dict(
            match={
                **where,
                cursorField: {
                    lt: before
                },
            },
        ))
    else:
        args = dict(match=where, )
    if pipeline:
        args.update(dict(pipeline=pipeline))
    args.update(dict(sort={cursorField: ASCENDING}))
    if first:
        args.update(dict(limit=first + 1, ))
    elif last:
        count = await mongodb_streams.count_documents(collection, where, pipeline=pipeline)
        toSkip = count - (last + 1)
        args.update(dict(skip=max(toSkip, 0)))

    nodes = await mongodb_streams.find(collection, **args)

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
            'endCursor': end_cursor and OUTPUT_COERCERS[scalar_name](end_cursor),
            'startCursor': start_cursor and OUTPUT_COERCERS[scalar_name](start_cursor),
            'hasNextPage': hasNext,
            'hasPreviousPage': hasPrevious,
        }
    }



MONGODB_OPERATORS = [
    'in',
    'nin',
    'eq',
    'neq',
    'or',
    'and',
    # TODO add gt, gte, like ....
]

def strip_nones(x: dict):
    result = {}
    for k, v in x.items():
        if not v == None:
            if k in MONGODB_OPERATORS:
                k = '$' + k
            if isinstance(v, dict):
                result[k] = strip_nones(v)
            else:
                result[k] = v
    return result

