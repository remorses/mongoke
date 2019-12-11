import json
import os.path
import os.path
import requests
from funcy import pluck, count, merge, collecting
from populate import indent_to

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

def download_file(url) -> str:
    r = requests.get(url, )
    return r.text

def get_types_schema(config, here='./'):
    if "schema_path" in config:
        path = here + config["schema_path"]
        with open(path) as f:
            return f.read()
    if config.get("schema"):
        return indent_to('',  config.get("schema"))
    if config.get("schema_url"):
        return download_file(config.get("schema_url"))
        

