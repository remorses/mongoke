from .support import zip_pluck, join_yields, repr_eval_dict
from populate import indent_to
import json
from funcy import lfilter, post_processing
from .resolvers_support import resolvers_support



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

@join_yields('')
def render_type_resolver(disambiguations, typename):
    code = f"""
    @TypeResolver('{typename}')
    def resolve_type(result, context, info, abstract_type):
        x = result
    """
    yield indent_to('', code) + '    '
    for (i, typename, expr) in zip_pluck(disambiguations, ['type_name', 'expression'], enumerate=True):
        code = f"""
        {'if' if i == 0 else 'elif'} ({expr}):
            return '{typename}'
        """ 
        yield indent_to('    ', code)


def repr_node_filterer(guards_after):
    code = f'''
    def filter_nodes_by_guard(nodes, fields, jwt):
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
    render_type_resolver=render_type_resolver,
)

resolvers_init = '''
from ..logger import logger
'''

generated_init = '''
from ..logger import logger
'''
# collection, resolver_path, guard_expression_before, guard_expression_after, disambiguations
single_item_resolver = '''
from tartiflette import Resolver, TypeResolver
from .support import strip_nones, zip_pluck
import mongodb_streams
from operator import setitem
from funcy import omit

${{render_type_resolver(disambiguations, typename) if disambiguations else ''}}

pipeline: list = ${{repr_eval_dict(pipeline,)}}

@Resolver('${{resolver_path}}')
async def resolve_${{'_'.join([x.lower() for x in resolver_path.split('.')])}}(parent, args, ctx, info):
    where = strip_nones(args.get('where', {}))
    headers = ctx['req'].headers
    jwt = ctx['req'].state.jwt_payload
    fields = []
    ${{repr_guards_checks(guards_before, '    ')}}
    collection = ctx['db']['${{collection}}']
    x = await mongodb_streams.find_one(collection, match=where, pipeline=pipeline)
    ${{repr_guards_checks(guards_after, '    ')}}
    # {{repr_disambiguations(disambiguations, '    ')}}
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
    jwt = ctx['req'].state.jwt_payload
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
        # direction=args
    )
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields, jwt=jwt))
    # {{repr_many_disambiguations(disambiguations, '    ') if disambiguations else ''}}
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
    # {{repr_disambiguations(disambiguations, '    ')}}
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
    jwt = ctx['req'].state.jwt_payload
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
    data['nodes'] = list(filter_nodes_by_guard(data['nodes'], fields, jwt=jwt))
    # {{repr_many_disambiguations(disambiguations, '    ') if disambiguations else ''}}
    return data
'''

