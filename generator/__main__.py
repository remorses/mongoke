import skema
from typing import *
import sys
from funcy import merge, lmap, collecting
import yaml
from skema.to_graphql import to_graphql
import os.path
from populate import populate_string
from .templates.resolvers import resolvers_init, resolvers_support, single_item_resolver, many_items_resolvers
from .templates.scalars import scalars_implementations
from .templates.graphql_query import graphql_query, general_graphql
from .templates.main import main
from .templates.jwt_middleware import jwt_middleware
from .templates.logger import logger
from .support import touch, pretty, zip_pluck, get_type_properties


def is_scalar(type_body):
    SCALARS = ['string', 'number', 'integer', 'boolean']
    return (
        type_body.get('type', '') in SCALARS
        # TODO add aliases, not only Type: Any
        or not {k: v for k, v in type_body.items() if k not in ('description', 'title',)}
    )


SCALAR_TYPES = ['String', 'Float', 'Int', 'Boolean', ]



@collecting
def get_scalar_fields(skema_schema, typename) -> List[Tuple[str, str]]:
    json_schema = skema.to_jsonschema(skema_schema, ref=typename, resolve=True)
    # pretty(json_schema)
    type_properties = get_type_properties(json_schema)
    aliases = skema.to_jsonschema(skema_schema, resolve=False).get('definitions', []).keys()
    aliases = list(aliases)
    pretty( aliases)
    for name, body in type_properties.items():
        if is_scalar(body):
            if body.get('title', ) in aliases: # TODO this logic is faulted, should be ported to skema, as get_schema_scalars
                scalar_name = body.get('title', )
            else:
                scalar_name = body.get('type', 'Unknown').capitalize()
            yield (name, scalar_name)


def add_guards_defaults(guard):
    guard['when'] = guard.get('when') or 'before'
    guard['fields'] = guard.get('fields') or []
    guard['roles'] = guard.get('roles') or []
    return guard


def add_disambiguations_defaults(dis):
    return dis


def generate_from_config(config):
    target_dir = config.get('target_dir', '.')
    root_dir_name = config.get('root_dir_name', 'root')
    base = os.path.join(target_dir, root_dir_name)
    skema_schema = config.get('skema')
    scalars = [*SCALAR_TYPES]  # TODO add other scalars from the skema
    main_graphql_schema = to_graphql(skema_schema)

    touch(f'{base}/__init__.py', '')
    touch(f'{base}/__main__.py', main)
    touch(f'{base}/generated/__init__.py', '')
    touch(f'{base}/generated/middleware/__init__.py', jwt_middleware)
    touch(f'{base}/generated/resolvers/__init__.py', resolvers_init)
    touch(f'{base}/generated/resolvers/support.py', resolvers_support)
    touch(f'{base}/generated/scalars/__init__.py', scalars_implementations)
    touch(f'{base}/generated/sdl/general.graphql',
          populate_string(general_graphql, dict(scalars=scalars)))
    touch(f'{base}/generated/sdl/main.graphql', main_graphql_schema)

    # needs:
    # disambiguations
    # typename
    # collection
    # guards

    for type_config in config.get('types', []):
        collection = type_config['collection']
        typename = type_config['type_name']
        guards = type_config.get('guards', [])
        guards = lmap(add_guards_defaults, guards)
        disambiguations = type_config.get('disambiguations', [])
        disambiguations = lmap(add_disambiguations_defaults, disambiguations)
        relations = type_config.get('relations', [])  # TODO relations

        query_name = typename[0].lower() + typename[1:]
        # pretty(get_scalar_fields(skema_schema, typename))
        query_subset = populate_string(
            graphql_query,
            dict(
                query_name=query_name,
                type_name=typename,
                fields=get_scalar_fields(skema_schema, typename),
                # scalars=scalars,
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
                disambiguations=disambiguations,
                guards_before=[g for g in guards if g['when'] == 'before'],
                guards_after=[g for g in guards if g['when'] == 'after'],
                zip_pluck=zip_pluck,
            )
        )
        touch(f'{base}/generated/resolvers/{query_name}.py', single_resolver)
        many_resolver = populate_string(
            many_items_resolvers,
            dict(
                # query_name=query_name,
                # type_name=typename,
                collection=collection,
                resolver_path='Query.' + query_name + 's',
                disambiguations=disambiguations,
                guards_before=[g for g in guards if g['when'] == 'before'],
                guards_after=[g for g in guards if g['when'] == 'after'],
                zip_pluck=zip_pluck,
            )
        )
        touch(f'{base}/generated/resolvers/{query_name}s.py', many_resolver)


config = yaml.safe_load(open(sys.argv[-1]).read())
generate_from_config(config)
