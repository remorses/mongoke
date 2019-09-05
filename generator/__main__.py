import skema
import requests
from typing import *
import sys
from funcy import merge, lmap, collecting, omit, remove
import yaml
from skema.to_graphql import to_graphql
import os.path
from populate import populate_string
from .templates.resolvers import resolvers_dependencies, resolvers_init, resolvers_support, single_item_resolver, many_items_resolvers, single_relation_resolver, many_relations_resolver
from .templates.scalars import scalars_implementations
from .templates.graphql_query import graphql_query, general_graphql, to_many_relation, to_many_relation_boilerplate, to_one_relation
from .templates.main import main
from .templates.jwt_middleware import jwt_middleware
from .templates.logger import logger
from .support import touch, pretty, get_type_properties

SCALAR_TYPES = ['String', 'Float', 'Int', 'Boolean', ]
SCALARS_ALREADY_IMPLEMENTED = ['ObjectId', 'Json',
                               'Date', 'DateTime', 'Time', *SCALAR_TYPES]


def is_scalar(type_body):
    SCALARS = ['string', 'number', 'integer', 'boolean']
    # print(omit(type_body, ['description', 'title', '$schema']))
    return (
        type_body.get('type', '') in SCALARS
        # TODO add aliases, not only Type: Any
        or not omit(type_body, ['description', 'title', '$schema'])
    )


def get_skema_aliases(skema_schema):
    definitions = skema.to_jsonschema(
        skema_schema, resolve=False).get('definitions', [])
    aliases = [body.get('title') for d, body in definitions.items()]
    # pretty(aliases)
    aliases = [x for x in aliases if is_alias(skema_schema, x)]
    return aliases


def is_alias(skema_schema, typename) -> bool:
    json_schema = skema.to_jsonschema(skema_schema, ref=typename, resolve=True)
    # pretty(json_schema)
    return is_scalar(json_schema)


@collecting
def get_scalar_fields(skema_schema, typename) -> List[Tuple[str, str]]:
    json_schema = skema.to_jsonschema(skema_schema, ref=typename, resolve=True)
    # pretty(json_schema)
    type_properties = get_type_properties(json_schema)
    aliases = get_skema_aliases(skema_schema)
    for name, body in type_properties.items():
        if is_scalar(body):
            #  TODO this logic is faulted, should be ported to skema, as get_schema_scalars
            if body.get('title', '') in aliases:
                scalar_name = body.get('title', )
            else:
                scalar_name = body.get('type', 'Unknown').capitalize()
            yield (name, scalar_name)


def add_guards_defaults(guard):
    guard['when'] = guard.get('when') or 'before'
    guard['fields'] = guard.get('fields') or []
    return guard


def add_disambiguations_defaults(dis):
    return dis


def generate_type_sdl(schema, typename, guards, query_name, is_aggregation=False):
    return populate_string(
        graphql_query,
        dict(
            query_name=query_name,
            type_name=typename,
            fields=get_scalar_fields(schema, typename),
            # scalars=scalars,
        )
    )



def generate_resolvers(collection, disambiguations, guards, query_name, is_aggregation=False, **kwargs):
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
            **resolvers_dependencies,
            **kwargs,
        )
    )
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
            **resolvers_dependencies,
            **kwargs,
        )
    )
    return single_resolver, many_resolver

@collecting
def make_disambiguations_objects(disambiguations):
    for type, expr in disambiguations.items():
        yield {
            'type_name': type,
            'expression': expr.strip(), 
        }

def get_skema(config):
    if 'skema_path' in config:
        with open(config['skema_path']) as f:
            return f.read()
    if config.get('skema'):
        return config.get('skema')
    if config.get('skema_url'):
        r = requests.get(config.get('skema_url'), stream=True)
        skema = []
        while 1:
            buf = r.raw.read(16*1024)
            if not buf:
                break
            skema.append(buf.decode())
        return skema
    

def generate_from_config(config):
    types = config.get('types', {})
    def get_type_config(name):
        if name in types.keys():
            conf = types[name]
        else:
            raise Exception(f'fromType {name} not found in config')
        return conf
    relations = config.get('relations', [])
    target_dir = config.get('target_dir', '.')
    root_dir_name = config.get('root_dir_name', 'root')
    base = os.path.join(target_dir, root_dir_name)
    skema_schema = get_skema(config)
    # TODO add other scalars from the skema
    scalars = [*SCALAR_TYPES, *get_skema_aliases(skema_schema)]
    main_graphql_schema = to_graphql(
        skema_schema, scalar_already_present=SCALARS_ALREADY_IMPLEMENTED)

    touch(f'{base}/__init__.py', '')
    touch(f'{base}/__main__.py', main)
    touch(f'{base}/generated/__init__.py', '')
    touch(f'{base}/generated/middleware/__init__.py', jwt_middleware)
    touch(f'{base}/generated/resolvers/__init__.py', resolvers_init)
    touch(f'{base}/generated/resolvers/support.py', resolvers_support)
    touch(f'{base}/generated/scalars/__init__.py',
          populate_string(scalars_implementations, dict(scalars=[x for x in scalars if x not in SCALARS_ALREADY_IMPLEMENTED])))
    touch(f'{base}/generated/sdl/general.graphql',
          populate_string(general_graphql, dict(scalars=scalars)))
    touch(f'{base}/generated/sdl/main.graphql', main_graphql_schema)

    # needs:
    # disambiguations
    # typename
    # collection
    # guards
    for typename, type_config in types.items():
        # types with no collection are used only for relations not direct queries
        if not type_config.get('exposed', True):
            continue
        collection = type_config['collection']
        query_name = typename[0].lower() + typename[1:]
        pipeline = type_config.get('pipeline', [])
        guards = type_config.get('guards', [])
        guards = lmap(add_guards_defaults, guards)
        disambiguations = type_config.get('disambiguations', {})
        disambiguations = make_disambiguations_objects(disambiguations)
        disambiguations = lmap(add_disambiguations_defaults, disambiguations)

        query_subset = generate_type_sdl(
            skema_schema,
            guards=guards,
            typename=typename,
            query_name=query_name
        )
        touch(f'{base}/generated/sdl/{query_name}.graphql', query_subset)
        single_resolver, many_resolver = generate_resolvers(
            collection=collection,
            disambiguations=disambiguations,
            guards=guards,
            query_name=query_name,
            pipeline=pipeline,
        )
        touch(f'{base}/generated/resolvers/{query_name}.py', single_resolver)
        touch(f'{base}/generated/resolvers/{query_name}s.py', many_resolver)

    for relation in relations:
        fromType = relation['from']
        fromTypeConfig = get_type_config(fromType)
        toType = relation['to']
        toTypeConfig = get_type_config(toType)
        relationName = relation.get('field')
        relation_type = relation.get('relation_type', 'to_one')
        relation_template = to_one_relation if relation_type == 'to_one' else to_many_relation
        relation_sdl = populate_string(
            relation_template,
            dict(
                toType=toType,
                fromType=fromType,
                relationName=relationName,
            )
        )
        implemented_types = [name for name, x in types.items() if x.get('exposed', True)]
        if relation_type == 'to_many' and toType not in implemented_types:
            relation_sdl += populate_string(
                to_many_relation_boilerplate,
                dict(
                    toType=toType,
                    fromType=fromType,
                    relationName=relationName,
                    fields=get_scalar_fields(skema_schema, toType),
                )
            )
        touch(f'{base}/generated/sdl/{fromType.lower()}_{relationName}.graphql', relation_sdl)
        relation_template = single_relation_resolver if relation_type == 'to_one' else many_relations_resolver
        relation_resolver = populate_string(
            relation_template,
            dict(
                # query_name=query_name,
                # type_name=typename,
                where_filter=relation['query'],
                pipeline=toTypeConfig.get('pipeline', []),
                collection=collection,
                resolver_path=fromType + '.' + relationName,
                # disambiguations=disambiguations,
                # guards_before=[g for g in guards if g['when'] == 'before'],
                # guards_after=[g for g in guards if g['when'] == 'after'],
                disambiguations=[],
                guards_before=[],
                guards_after=[],
                **resolvers_dependencies,
            )
        )
        touch(f'{base}/generated/resolvers/{fromType.lower()}_{relationName}.py', relation_resolver)


arg = sys.argv[-1]
arg = 'pr_conf.yaml'
config = yaml.safe_load(open(arg).read())
generate_from_config(config)

