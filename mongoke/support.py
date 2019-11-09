import json
import os.path
import os.path
import requests
from funcy import pluck, count, merge, collecting
from populate import indent_to
from skema.reconstruct import from_jsonschema
from skema import gen_graphql
from .logger import logger

def get_config_schema():
    with open(os.path.dirname(__file__) +  '/config_schema.json') as f:
        return json.loads(f.read())


def to_string(i, length=2):
    i = str(i)
    prefix = '0' * (length - len(i))
    return prefix + i

i = 0
def make_touch(base):
    def touch(filename, data, index=False):
        global i
        if index:
            parts = filename.split('/')
            parts = parts[:-1] + [f'{to_string(i)}_' + parts[-1]]
            filename = '/'.join(parts)
            i += 1
        filename = f'{base}/{filename}'
        os.makedirs(os.path.dirname(filename), exist_ok=True)
        with open(filename, "w") as f:
            f.write(data)
    return touch


pretty = lambda x: print(json.dumps(x, indent=4, default=str))


def zip_pluck(d, keys, enumerate=False):
    args = [pluck(k, d) for k in keys]
    if enumerate:
        args = [count(), *args]
    return zip(*args)

def find(l, predicate):
    for x in l:
        if (predicate(x)):
            return x
    return None

@collecting
def unique(l, key=lambda x: x):
    found = []
    for x in l:
        id = key(x)
        if not id in found:
            found += [id]
            yield x

def read(path):
        with open(path) as f:
            return f.read()

def jsonschema_to_graphql(schema):
        tree = from_jsonschema(schema)
        logger.info('generating graphql from jsonschema')
        graphql = gen_graphql(tree)
        logger.info(graphql)
        return graphql

def get_types_schema(config, here='./'):
    if config.get("jsonschema"):
        schema = indent_to('',  config.get("jsonschema"))
        return jsonschema_to_graphql(schema)
    if config.get("schema"):
        return indent_to('',  config.get("schema"))
    if "jsonschema_path" in config:
        path = here + config["jsonschema_path"]
        schema = json.loads(read(path))
        return jsonschema_to_graphql(schema)
    if "schema_path" in config:
        path = here + config["schema_path"]
        return read(path)

    if config.get("schema_url"):
        r = requests.get(config.get("schema_url"), stream=True)
        skema = ""
        while 1:
            buf = r.raw.read(16 * 1024)
            if not buf:
                break
            skema += buf.decode()
        return skema

