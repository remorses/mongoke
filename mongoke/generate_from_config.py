import skema
from shutil import rmtree
import requests
from typing import *
import sys
from funcy import merge, lmap, collecting, omit, remove, lcat
from .checksum import make_checksum, existent_checksum
import yaml
from skema.to_graphql import to_graphql
import os.path
from populate import populate_string
from .templates.resolvers import (
    resolvers_dependencies,
    resolvers_init,
    resolvers_support,
    single_item_resolver,
    many_items_resolvers,
    single_relation_resolver,
    many_relations_resolver,
    generated_init,
)
from .templates.scalars import scalars_implementations
from .templates.graphql_query import (
    graphql_query,
    general_graphql,
    to_many_relation,
    to_many_relation_boilerplate,
    to_one_relation,
)
from .templates.main import main
from .templates.jwt_middleware import jwt_middleware
from .templates.logger import logger
from .templates.engine import engine
from .support import make_touch, pretty, get_skema
from .skema_support import get_scalar_fields, get_skema_aliases
from .naming import get_query_name, get_relation_filename, get_resolver_filenames
from .generators import generate_relation_boilerplate, generate_type_boilerplate

SCALAR_TYPES = ["String", "Float", "Int", "Boolean", "Json"]
SCALARS_ALREADY_IMPLEMENTED = [
    "ID",
    "ObjectId",
    "Json",
    "Date",
    "DateTime",
    "Time",
    *SCALAR_TYPES,
]


def add_guards_defaults(guard):
    guard["when"] = guard.get("when") or "after"
    guard["excluded"] = guard.get("excluded") or []
    return guard


def add_disambiguations_defaults(dis):
    return dis


@collecting
def make_disambiguations_objects(disambiguations):
    for type, expr in disambiguations.items():
        yield {"type_name": type, "expression": expr.strip()}


def generate_from_config(config, config_path, root_dir_path):
    types = config.get("types", {})
    jwt_config = config.get("jwt", {})
    relations = config.get("relations", [])
    root_dir_path = root_dir_path or "generated"
    if os.path.exists(root_dir_path):
        rmtree(root_dir_path)
    db_url = config.get("db_url", "")
    touch = make_touch(base=os.path.abspath(root_dir_path))
    skema_schema = get_skema(config, here=config_path)

    # TODO add other scalars from the skema
    scalars = {*SCALAR_TYPES, *get_skema_aliases(skema_schema)}
    scalars = list(scalars)
    main_graphql_schema = to_graphql(
        skema_schema, hide=SCALARS_ALREADY_IMPLEMENTED
    )

    touch(f"checksum", make_checksum(config, config_path))
    touch(f"__init__.py", "")
    touch(f"engine.py", engine)
    touch(
        f"__main__.py",
        populate_string(
            main,
            dict(
                root_dir_name=root_dir_path.split("/")[-1],
                db_url=db_url,
                resolver_names=get_resolver_filenames(config),
            ),
        ),
    )
    touch(f"generated/__init__.py", "")
    touch(f"generated/logger.py", logger)
    touch(f"generated/middleware/__init__.py", 
        populate_string(
            jwt_middleware,
            dict(
                jwt_header=jwt_config.get('header_name',) or 'Authorization',
                jwt_sheme=jwt_config.get('header_scheme', 'Bearer'),
                jwt_required=bool(jwt_config.get('required')),
                jwt_secret=jwt_config.get('secret', None),
                jwt_algorithms=jwt_config.get('algorithms', ['H256']),
            ),
        ),
    )

    touch(f"generated/resolvers/__init__.py", resolvers_init)
    touch(
        f"generated/resolvers/support.py",
        populate_string(
            resolvers_support,
            dict(scalars=[x for x in scalars if x not in SCALARS_ALREADY_IMPLEMENTED]),
        ),
    )
    touch(
        f"generated/scalars.py",
        populate_string(
            scalars_implementations,
            dict(scalars=[x for x in scalars if x not in SCALARS_ALREADY_IMPLEMENTED]),
        ),
    )
    touch(
        f"generated/sdl/general.graphql",
        populate_string(general_graphql, dict(scalars=scalars)),
        index=True
    )
    touch(f"generated/sdl/main.graphql", main_graphql_schema, index=True)
    implemented_types = []
    for typename, type_config in types.items():
        type_config = type_config or {}
        if not type_config.get("exposed", True):
            continue
        disambiguations = type_config.get("disambiguations", {})
        disambiguations = make_disambiguations_objects(disambiguations)
        disambiguations = lmap(add_disambiguations_defaults, disambiguations)
        generate_type_boilerplate(
            touch=touch,
            skema_schema=skema_schema,
            collection=type_config.get("collection", ""),
            typename=typename,
            guards=lmap(add_guards_defaults, type_config.get("guards", [])),
            pipeline=type_config.get("pipeline", []),
            disambiguations=disambiguations,
        )
        implemented_types += [typename]

    for relation in relations:
        toType = relation["to"]
        generate_relation_boilerplate(
            touch=touch,
            skema_schema=skema_schema,
            fromType=relation["from"],
            where_filter=relation["where"],
            toType=toType,
            pipeline=types[toType].get("pipeline", []),
            collection=types[toType].get("collection", []),
            relationName=relation.get("field"),
            relation_type=relation.get("relation_type", "to_one"),
            implemented_types=implemented_types,
            resolver_filename=get_relation_filename(relation),
        )
        implemented_types += [toType]

