from funcy import collecting


def get_relation_filename(relation):
    fromType = relation["from"]
    relationName = relation["field"]
    return f"{fromType}_{relationName}"


@collecting
def get_resolver_filenames(config):
    for typename, type_config in config.get("types", {}).items():
        if type_config.get("exposed", True):
            yield get_query_name(typename)
            yield get_query_name(typename) + "s"
    for relation in config.get("relations", []):
        yield get_relation_filename(relation)


def get_query_name(typename):
    return typename
    # return typename[0].lower() + typename[1:]
