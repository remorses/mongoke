from .support import zip_pluck, join_yields, repr_eval_dict
from populate import indent_to
import json
from funcy import lfilter, post_processing



@join_yields('')
def repr_guards_checks(guards, indentation):
    for expr, fields in zip_pluck(guards, ['expression', 'excluded']):
        code =  f"""
        if not ({expr}):
            raise Exception({json.dumps('guard `' + str(expr) + '` not satisfied')})
        else:
            fields += {fields}
        """
        yield indent_to(indentation, code)


@join_yields('')
def repr_disambiguations(disambiguations, indentation):
    for (i, typename, expr) in zip_pluck(disambiguations, ['type_name', 'expression'], enumerate=True):
        code = f"""
        {'if' if i == 0 else 'elif'} ({expr}):
            x['_typename'] = '{typename}'
        """ 
        yield indent_to(indentation, code)


def repr_node_filterer(guards_after):
    code = f'''
    def filter_nodes_by_guard(nodes, fields):
        for x in nodes:
            try:
                {repr_guards_checks(guards_after, '                ')}
                yield omit(x or dict(), fields)
            except Exception:
                pass
    '''
    return indent_to('', code)

def repr_many_disambiguations(disambiguations, indentation):
    code = f'''
    for x in data['nodes']:
        {repr_disambiguations(disambiguations, '        ')}
    '''
    return indent_to(indentation, code)

resolvers_dependencies = dict(
    repr_guards_checks=repr_guards_checks,
    zip_pluck=zip_pluck,
    repr_disambiguations=repr_disambiguations,
    repr_eval_dict=repr_eval_dict,
    repr_node_filterer=repr_node_filterer,
    repr_many_disambiguations=repr_many_disambiguations,
)

resolvers_init = '''
from ..logger import logger
'''

generated_init = '''
from ..logger import logger
'''
# collection, resolver_path, guard_expression_before, guard_expression_after, disambiguations
single_item_resolver = '''
from tartiflette import Resolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem
from funcy import omit

pipeline: list = ${{repr_eval_dict(pipeline,)}}

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    ${{repr_guards_checks(guards_before, '    ')}}
    collection = ctx['db']['${{collection}}']
    x = await mongodb_streams.find_one(collection, match=where, pipeline=pipeline)
    ${{repr_guards_checks(guards_after, '    ')}}
    ${{repr_disambiguations(disambiguations, '    ')}}
    if fields:
        x = omit(x or dict(), fields)
    return x
'''

# collection, resolver_path, guard_expression_before, guard_expression_after, disambiguations
many_items_resolvers = '''
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem
from funcy import omit

${{repr_node_filterer(guards_after)}}

map_fields_to_types = ${{repr_eval_dict(map_fields_to_types, '    ')}}

pipeline: list = ${{repr_eval_dict(pipeline,)}}

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    cursorField = args.get('cursorField',) or ('_id' if '_id' in map_fields_to_types else list(map_fields_to_types.keys())[0])
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload
    fields = []
    ${{repr_guards_checks(guards_before, '    ')}}
    pagination = get_pagination(args,)
    data = await connection_resolver(
        collection=ctx['db']['${{collection}}'], 
        where=where,
        cursorField=cursorField,
        pagination=pagination,
        scalar_name=map_fields_to_types[cursorField],
        pipeline=pipeline,
    )
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields))
    ${{repr_many_disambiguations(disambiguations, '    ') if disambiguations else ''}}
    return data

'''

# where_filter, collection, resolver_path
# TODO add guards, disambig
# TODO add pipeline for making an aggregate
single_relation_resolver = ''' 
from tartiflette import Resolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem

pipeline: list = ${{repr_eval_dict(pipeline,)}}

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = ${{repr_eval_dict(where_filter, '    ')}}
    ${{repr_guards_checks(guards_before, '    ')}}
    collection = ctx['db']['${{collection}}']
    x = await mongodb_streams.find_one(collection, match=where, pipeline=pipeline)
    ${{repr_guards_checks(guards_after, '    ')}}
    ${{repr_disambiguations(disambiguations, '    ')}}
    return x
'''

# where_filter, collection
# TODO add pipeline for making an aggregate
many_relations_resolver = '''
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem
from funcy import omit

${{repr_node_filterer(guards_after)}}

map_fields_to_types = ${{repr_eval_dict(map_fields_to_types, '    ')}}

pipeline: list = ${{repr_eval_dict(pipeline,)}}

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    relation_where = ${{repr_eval_dict(where_filter, '    ')}}
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)
    cursorField = args.get('cursorField',) or ('_id' if '_id' in map_fields_to_types else list(map_fields_to_types.keys())[0])
    headers = ctx['req'].headers
    jwt = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
    ${{repr_guards_checks(guards_before, '    ')}}
    pagination = get_pagination(args,)
    data = await connection_resolver(
        collection=ctx['db']['${{collection}}'], 
        where=where,
        cursorField=cursorField,
        pagination=pagination,
        scalar_name=map_fields_to_types[cursorField],
        pipeline=pipeline,
    )
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields))
    ${{repr_many_disambiguations(disambiguations, '    ') if disambiguations else ''}}
    return data
'''

# nothing
resolvers_support = '''
import collections
from prtty import pretty
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
    'ID': str,
    ${{'**{scalar.name: scalar._implementation.coerce_input for scalar in scalar_classes},' }}
}

OUTPUT_COERCERS = {
    'String': str,
    'Int': int,
    'Float': float,
    'Bool': bool,
    'ID': str,
    ${{'**{scalar.name: scalar._implementation.coerce_output for scalar in scalar_classes},' }}
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
    pretty(args)
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

'''