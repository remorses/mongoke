
import collections
import os
from prtty import pretty
from motor.motor_asyncio import AsyncIOMotorDatabase, AsyncIOMotorCollection
import mongodb_streams
from tartiflette import Resolver
import pymongo
from pymongo import ASCENDING, DESCENDING
from typing import NamedTuple, Union
import typing
from funcy import pluck, select_keys, omit, lmap
from ..scalars import scalar_classes



DEFAULT_NODES_COUNT = 20

INPUT_COERCERS = {
    None: lambda x: x,
    'String': str,
    'Int': int,
    'Float': float,
    'Bool': bool,
    'ID': str,
    **{scalar.name: scalar._implementation.coerce_input for scalar in scalar_classes},
}

OUTPUT_COERCERS = {
    None: lambda x: x,
    'String': str,
    'Int': int,
    'Float': float,
    'Bool': bool,
    'ID': str,
    **{scalar.name: scalar._implementation.coerce_output for scalar in scalar_classes},
}

def zip_pluck(d, *keys):
    return zip(*[pluck(k, d) for k in keys])

direction_map = {
    'ASC': ASCENDING,
    'DESC': DESCENDING,
}

def get_pagination(args,):
    after = args.get('after')
    before = args.get('before')
    return {
        'after': after,
        'before': before,
        'first': args.get('first'),
        'last': args.get('last'),
        'direction': direction_map[args.get('direction', 'DESC')],
    }

def opposite_direction(dir):
    if dir == ASCENDING:
        return DESCENDING
    return ASCENDING

async def connection_resolver(
    collection: AsyncIOMotorCollection,
    where: dict,
    cursorField,  # needs to exist always at least one, the fisrst is the cursorField
    pagination: dict,
    scalar_name,
    pipeline=[],
):
    if os.getenv('DEBUG'):
        print('executing connection_resolver')
        pretty({
            'where': where,
            'cursorField': cursorField,
            'pagination': pagination,
            'scalar_name': scalar_name,
            'collection': collection,
            'pipeline': pipeline,
        })
    direction = pagination['direction']
    first, last = pagination.get('first'), pagination.get('last'),
    after, before = pagination.get('after'), pagination.get('before')
    if after:
        after = INPUT_COERCERS.get(scalar_name, lambda x: x)(after)
    if before:
        before = INPUT_COERCERS.get(scalar_name, lambda x: x)(before)

    first = first or 0
    last = last or 0

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
    lt = '$gt' if direction == DESCENDING else '$lt'
    gt = '$lt' if direction == DESCENDING else '$gt'
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
    sorting = direction if not last else opposite_direction(direction)
    args.update(dict(sort={cursorField: sorting}))
    if last:
        args.update(dict(limit=last + 1, ))
    if first:
        args.update(dict(limit=first + 1, ))
    # elif first:
    #     count = await mongodb_streams.count_documents(collection, args['match'], pipeline=pipeline)
    #     toSkip = count - (last + 1)
    #     args.update(dict(skip=max(toSkip, 0)))
    args.update(dict(max_len=10000))
    # pretty(args)
    nodes = await mongodb_streams.find(collection, **args)

    hasNext = None
    hasPrevious = None

    if first:
        hasNext = len(nodes) == (first + 1)
        nodes = nodes[:-1] if hasNext else nodes

    if last:
        nodes = list(reversed(nodes))
        hasPrevious = len(nodes) == (last + 1)
        nodes = nodes[1:] if hasPrevious else nodes

    end_cursor = nodes[-1].get(cursorField) if nodes else None
    start_cursor = nodes[0].get(cursorField) if nodes else None
    return {
        'nodes': nodes,
        'edges': lmap(
            lambda node: dict(
                node=node, 
                cursor=OUTPUT_COERCERS[scalar_name](node.get(cursorField))
            ), nodes),
        'pageInfo': {
            'endCursor': end_cursor and OUTPUT_COERCERS[scalar_name](end_cursor),
            'startCursor': start_cursor and OUTPUT_COERCERS[scalar_name](start_cursor),
            'hasNextPage': hasNext,
            'hasPreviousPage': hasPrevious,
        }
    }

def make_edge(node, cursorField):
    return {
        'node': node,
        'cursor': node.get(cursorField),
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
        if not v == None and v != {}:
            if k in MONGODB_OPERATORS:
                k = '$' + k
            if isinstance(v, dict):
                result[k] = strip_nones(v)
            else:
                result[k] = v
    return result

