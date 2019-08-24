import skema
import sys
from funcy import merge
import yaml
from skema.to_graphql import to_graphql
import os.path
from populate import populate_string
from .templates.resolvers import resolvers_init, resolvers_support, single_item_resolver, many_items_resolvers
from .templates.scalars import scalars
from .templates.graphql_query import graphql_query, general_graphql
from .templates.main import main
from .templates.jwt_middleware import jwt_middleware
from .templates.logger import logger
from .support import touch, pretty

def is_scalar(type_body):
    SCALARS = ['string', 'number', 'integer', 'boolean']
    return (
        type_body.get('type', '') in SCALARS 
        or not {k:v for k, v in type_body.items() if k not in ('description', 'title',)}
    )


def get_scalar_fields(skema_schema, typename):
    json_schema = skema.to_jsonschema(skema_schema, ref=typename, resolve=True)
    # pretty(json_schema)
    if any([x in json_schema for x in ('anyOf', 'allOf', 'oneOf')]):
        subsets = json_schema.get('anyOf', [])
        subsets = subsets or json_schema.get('allOf', [])
        subsets = subsets or json_schema.get('oneOf', [])
        type_properties = merge(*[x.get('properties',) for x in subsets])
    else:
        type_properties = json_schema.get('properties', {})
    return [name for name, body in type_properties.items() if is_scalar(body)]


def generate_from_config(config):
    target_dir = config.get('target_dir', '.')
    root_dir_name = config.get('root_dir_name', 'root')
    base = os.path.join(target_dir, root_dir_name)
    skema_schema = config.get('skema')
    main_graphql_schema = to_graphql(skema_schema)

    touch(f'{base}/__init__.py', '')
    touch(f'{base}/__main__.py', main)
    touch(f'{base}/generated/__init__.py', '')
    touch(f'{base}/generated/middleware/__init__.py', jwt_middleware)
    touch(f'{base}/generated/resolvers/__init__.py', resolvers_init)
    touch(f'{base}/generated/resolvers/support.py', resolvers_support)
    touch(f'{base}/generated/scalars/__init__.py', scalars)
    touch(f'{base}/generated/sdl/general.graphql', general_graphql)
    touch(f'{base}/generated/sdl/main.graphql', main_graphql_schema)


    # needs:
        # disambiguations
        # typename
        # collection
        # guards

    for typename, collection in config.get('database', {}).get('collections', {}).items():
        if collection:
            query_name = typename[0].lower() + typename[1:]
            # pretty(get_scalar_fields(skema_schema, typename))
            query_subset = populate_string(
                graphql_query,
                dict(
                    query_name=query_name,
                    type_name=typename,
                    fields=get_scalar_fields(skema_schema, typename),
                )
            )
            touch(f'{base}/generated/sdl/{query_name}.graphql', query_subset)

            single_resolver = populate_string(
                single_item_resolver,
                dict(
                    # query_name=query_name,
                    # type_name=typename,
                    collection=collection,
                    resolver_path='Query.' + query_name,
                    disambiguations={},
                    guard_expression_before=True,
                    guard_expression_after=True,
                )
            )
            touch(f'{base}/generated/resolvers/{query_name}.py', single_resolver)



config = yaml.safe_load(open(sys.argv[-1]).read())
generate_from_config(config)


