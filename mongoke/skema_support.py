import skema
from typing import *
from funcy import merge, lmap, collecting, omit, remove, lcat
import yaml
from skema.to_graphql import to_graphql


HIDE_GRAPHQL = '[graphql hide]'


def get_skema_aliases(skema_schema):
    definitions = skema.to_jsonschema(
        skema_schema, resolve=False).get('definitions', [])
    definitions = {d: body for d, body in definitions.items() if not HIDE_GRAPHQL in body.get(
        'description', '') and not body.get('enum')}  # TODO should be implemented in skema
    aliases = [body.get('title') for d, body in definitions.items()]
    # pretty(aliases)
    aliases = [x for x in aliases if is_alias(skema_schema, x)]
    return aliases


@collecting
def get_scalar_fields(skema_schema, typename) -> Iterable[Tuple[str, str]]:
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
                scalar_name = map_json_type_to_grpahql[body.get(
                    'type', None).lower()]
            yield (name, scalar_name)


def get_type_properties(json_schema):
    if any([x in json_schema for x in ('anyOf', 'allOf', 'oneOf')]):
        subsets = json_schema.get('anyOf', [])
        subsets = subsets or json_schema.get('allOf', [])
        subsets = subsets or json_schema.get('oneOf', [])
        type_properties = merge(*[x.get('properties',) for x in subsets])
    else:
        type_properties = json_schema.get('properties', {})
    return type_properties


map_json_type_to_grpahql = {
    'string': 'String',
    'number': 'Float',
    'integer': 'Int',
    'boolean': 'Boolean',
    None: 'Json',
}


def is_scalar(type_body):
    SCALARS = ['string', 'number', 'integer', 'boolean']
    # print(omit(type_body, ['description', 'title', '$schema']))
    return (
        type_body.get('type', '') in SCALARS
        # TODO add aliases, not only Type: Any
        or not omit(type_body, ['description', 'title', '$schema'])
    )


def is_alias(skema_schema, typename) -> bool:
    json_schema = skema.to_jsonschema(skema_schema, ref=typename, resolve=True)
    # pretty(json_schema)
    return is_scalar(json_schema)
