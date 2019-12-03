from typing import *
from funcy import merge, lmap, collecting, omit, remove, lcat
from functools import lru_cache
import yaml
from graphql import parse, DocumentNode, Node, ListTypeNode
from graphql.language import (
    UnionTypeDefinitionNode,
    ScalarTypeDefinitionNode,
    FragmentDefinitionNode,
    EnumTypeDefinitionNode,
)
from .support import find, unique
from .constants import SCALAR_TYPES


@lru_cache(maxsize=10)
def parse_graphql_schema(schema) -> DocumentNode:
    doc = parse(schema)
    return doc


def get_fields(graphql_schema, typename) -> Iterable[Tuple[str, str]]:
    IGNORE_FIELDS = (
        ScalarTypeDefinitionNode,
        FragmentDefinitionNode,
        EnumTypeDefinitionNode,
    )
    doc = parse_graphql_schema(graphql_schema)
    node: Node = find(doc.definitions, lambda x: x.name.value == typename)
    if getattr(node, "fields", None):
        return [(field.name.value, get_type_name(field.type)) for field in node.fields]
    else:
        if isinstance(node, (UnionTypeDefinitionNode,)):
            fields = [get_fields(graphql_schema, x.name.value) for x in node.types]
            fields = unique(lcat(fields), key=lambda x: x[0])
            return fields
        elif isinstance(node, IGNORE_FIELDS):
            print(f"ignoring {node}")
            return []
        else:
            raise Exception(f"unrecognized type for {node}")

    # print(f"can't get fields, {e}")
    # return []


def get_type_name(node: Node):
    try:
        return node.name.value
    except Exception as e:
        if isinstance(node, ListTypeNode):
            return "LISTA"
        print(f"can't get type name, {e}")
        return "NOT_FOUND"


def get_scalar_fields(graphql_schema, typename):
    fields = get_fields(graphql_schema, typename)
    fields = [
        (name, _type) for (name, _type) in fields if is_scalar(graphql_schema, _type)
    ]
    return fields


def is_scalar(schema, typename):
    ok = typename in SCALAR_TYPES
    ok = ok or typename in get_graphql_scalars(schema)
    ok = ok or typename in get_graphql_enums(schema)
    return ok


@collecting
def get_graphql_types(schema, instanceof) -> List[str]:  # TODO
    doc = parse_graphql_schema(schema)
    for node in doc.definitions:
        if isinstance(node, instanceof):
            yield node.name.value


get_graphql_scalars = lambda schema: get_graphql_types(
    schema, instanceof=ScalarTypeDefinitionNode
)

get_graphql_enums = lambda schema: get_graphql_types(
    schema, instanceof=EnumTypeDefinitionNode
)

