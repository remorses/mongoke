from .support import zip_pluck, indent_to, join_yields, repr_eval_dict
import json
from funcy import lfilter, post_processing



@join_yields('\n')
def repr_guards_before_checks(guards_before, indentation):
    for expr, fields in zip_pluck(guards_before, ['expression', 'fields']):
        code =  f"""
        if not ({expr}):
            raise Exception({json.dumps('guard `' + str(expr) + '` not satisfied')})
        else:
            fields += {fields}
        """
        yield indent_to(indentation, code)

@join_yields('\n')
def repr_guards_after_checks(guards_after, indentation):
    for expr, fields in zip_pluck(guards_after, ['expression', 'fields']):
        code = f"""
        if not ({expr}):
            raise Exception({json.dumps('guard `' + str(expr) + '` not satisfied')})
        else:
            fields += {fields}
        """
        yield indent_to(indentation, code)


@join_yields('\n')
def repr_disambiguations(disambiguations, indentation):
    for (i, typename, expr) in zip_pluck(disambiguations, ['type_name', 'expression'], enumerate=True):
        code = f"""
        {'if' if i == 0 else 'elif'} ({expr}):
            x['_typename'] = '{typename}'
        """ 
        yield indent_to(indentation, code)

@join_yields('\n')
def filter_nodes_guards_after(guards_after, indentation):
    for expr, fields in zip_pluck(guards_after, ['expression', 'fields']):
        code = f"""
        if not ({expr}):
            pass
        else:
            own_fields = fields + {fields}
            if own_fields:
                x = select_keys(lambda k: k in fields, x)
            nodes.append(x)
        """
        yield indent_to(indentation, code)


resolvers_dependencies = dict(
    repr_guards_before_checks=repr_guards_before_checks,
    repr_guards_after_checks=repr_guards_after_checks,
    filter_nodes_guards_after=filter_nodes_guards_after,
    zip_pluck=zip_pluck,
    repr_disambiguations=repr_disambiguations,
    repr_eval_dict=repr_eval_dict,
)

resolvers_init = '''
from ..logger import logger
'''
# collection, resolver_path, guard_expression_before, guard_expression_after, disambiguations
single_item_resolver = '''
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys
from operator import setitem
from funcy import select_keys

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['request']['headers']
    jwt = ctx['req'].jwt_payload # TODO i need to decode jwt_payload and set it in req in a middleware
    fields = []
${{repr_guards_before_checks(guards_before, '    ')}}
    collection = ctx['db']['${{collection}}']
    x = collection.find_one(where)
${{repr_guards_after_checks(guards_after, '    ')}}
${{repr_disambiguations(disambiguations, '    ')}}
    if fields:
        x = select_keys(lambda k: k in fields, x)
    return x
'''

# collection, resolver_path, guard_expression_before, guard_expression_after, disambiguations
many_items_resolvers = '''
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    orderBy = args.get('orderBy', {'_id': 'ASC'}) # add default
    headers = ctx['request']['headers']
    jwt = ctx['req'].jwt_payload # TODO i need to decode jwt_payload
    fields = []
${{repr_guards_before_checks(guards_before, '    ')}}
    pagination = get_pagination(args)
    data = await connection_resolver(
        collection=ctx['db']['${{collection}}'], 
        where=where,
        orderBy=orderBy,
        pagination=pagination,
    )
${{
"""
    nodes = []
    for x in data['nodes']:
""" if guards_after else """
    nodes = data['nodes']
""" if disambiguations else ''
}}
${{filter_nodes_guards_after(guards_after, '        ')}}
${{
"""
    for x in nodes:
""" if disambiguations else ''
}}
${{repr_disambiguations(disambiguations, '        ')}}
    ${{"data['nodes'] = nodes" if guards_after or disambiguations else ''}}
    return data

'''

# where_filter, collection, resolver_path
# TODO add guards, disambig
# TODO add pipeline for making an aggregate
single_relation_resolver = ''' 
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys
from operator import setitem

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = ${{repr_eval_dict(where_filter, '    ')}}
${{repr_guards_before_checks(guards_before, '    ')}}
    x = await ctx['db']['${{collection}}'].find_one(where)
${{repr_guards_after_checks(guards_after, '    ')}}
${{repr_disambiguations(disambiguations, '    ')}}
    return x
'''

# where_filter, collection
# TODO add pipeline for making an aggregate
many_relations_resolver = '''
from tartiflette import Resolver
from .support import strip_nones, connection_resolver, zip_pluck, select_keys, get_pagination
from operator import setitem

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    relation_where = ${{repr_eval_dict(where_filter, '    ')}}
    where = {**args.get('where', {}), **relation_where}
    where = strip_nones(where)
    
'''

# nothing
resolvers_support = '''
import collections
from motor.motor_asyncio import AsyncIOMotorDatabase, AsyncIOMotorCollection
from tartiflette import Resolver
import pymongo
from pymongo import ASCENDING, DESCENDING
from typing import NamedTuple, Union
import typing
from funcy import pluck, select_keys

gt = '$gt'
lt = '$lt'
MAX_NODES = 20

def zip_pluck(d, *keys):
    return zip(*[pluck(k, d) for k in keys])

def get_pagination(args):
    return {
        'after': args.get('after'),
        'before': args.get('before'),
        'first': args.get('first'),
        'last': args.get('last'),
    }


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

'''