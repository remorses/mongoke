import os.path
import sys
from typing import *

import requests

import yaml
from funcy import collecting, lcat, lmap, merge, omit, remove
from populate import populate_string

from .naming import get_query_name, get_relation_filename, get_resolver_filenames
from .graphql_support import get_scalar_fields
from .support import get_types_schema, make_touch, pretty
from .templates.engine import engine
from .templates.graphql_query import (
    general_graphql,
    graphql_query,
    to_many_relation,
    to_many_relation_boilerplate,
    to_one_relation,
)
from .templates.jwt_middleware import jwt_middleware
from .templates.logger import logger

from .templates.resolvers import (
    generated_init,
    many_items_resolvers,
    many_relations_resolver,
    resolvers_dependencies,
    resolvers_init,
    resolvers_support,
    single_item_resolver,
    single_relation_resolver,
)
from .templates.scalars import scalars_implementations


def generate_resolvers(
    collection, disambiguations, guards, typename, query_name, is_aggregation=False, **kwargs
):
    single_resolver = populate_string(
        single_item_resolver,
        dict(
            #  query_name=query_name,
            # type_name=typename,
            typename=typename,
            collection=collection,
            resolver_path="Query." + query_name,
            disambiguations=disambiguations,
            guards_before=[g for g in guards if g["when"] == "before"],
            guards_after=[g for g in guards if g["when"] == "after"],
            **resolvers_dependencies,
            **kwargs,
        ),
    )
    many_resolver = populate_string(
        many_items_resolvers,
        dict(
            #  query_name=query_name,
            typename=typename,
            collection=collection,
            resolver_path="Query." + query_name + "s",
            disambiguations=disambiguations,
            guards_before=[g for g in guards if g["when"] == "before"],
            guards_after=[g for g in guards if g["when"] == "after"],
            **resolvers_dependencies,
            **kwargs,
        ),
    )
    return single_resolver, many_resolver


def generate_type_sdl(schema, typename, guards, query_name, is_aggregation=False):
    return populate_string(
        graphql_query,
        dict(
            query_name=query_name,
            type_name=typename,
            fields=get_scalar_fields(schema, typename),
            # scalars=scalars,
        ),
    )


def generate_type_boilerplate(
    touch, schema, typename, guards, disambiguations, collection, pipeline
):
    query_name = get_query_name(typename)

    query_subset = generate_type_sdl(
        schema, guards=guards, typename=typename, query_name=query_name
    )
    touch(f"generated/sdl/{query_name}.graphql", query_subset, index=True)
    single_resolver, many_resolver = generate_resolvers(
        typename=typename,
        collection=collection,
        disambiguations=disambiguations,
        guards=guards,
        query_name=query_name,
        pipeline=pipeline,
        map_fields_to_types=dict(get_scalar_fields(schema, typename)),
    )
    touch(f"generated/resolvers/{get_query_name(typename)}.py", single_resolver)
    touch(f"generated/resolvers/{get_query_name(typename)}s.py", many_resolver)


def generate_relation_boilerplate(
    where_filter,
    touch,
    schema,
    relation_type,
    relationName,
    fromType,
    toType,
    implemented_types,
    pipeline,
    collection,
    resolver_filename,
):
    relation_template = (
        to_one_relation if relation_type == "to_one" else to_many_relation
    )
    relation_sdl = populate_string(
        relation_template,
        dict(toType=toType, fromType=fromType, relationName=relationName),
    )
    if relation_type == "to_many" and toType not in implemented_types:
        relation_sdl += populate_string(
            to_many_relation_boilerplate,
            dict(
                toType=toType,
                fromType=fromType,
                relationName=relationName,
                fields=get_scalar_fields(schema, toType),
            ),
        )
    touch(f"generated/sdl/{fromType}_{relationName}.graphql", relation_sdl, index=True)
    relation_template = (
        single_relation_resolver
        if relation_type == "to_one"
        else many_relations_resolver
    )
    relation_resolver = populate_string(
        relation_template,
        dict(
            where_filter=where_filter,
            pipeline=pipeline,
            collection=collection,
            resolver_path=fromType + "." + relationName,
            disambiguations=[],  # TODO relation guards
            guards_before=[],
            guards_after=[],
            map_fields_to_types=dict(get_scalar_fields(schema, toType)),
            **resolvers_dependencies,
        ),
    )
    touch(f"generated/resolvers/{resolver_filename}.py", relation_resolver)
