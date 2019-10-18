import json
import os.path
import os.path
import requests
from funcy import pluck, count, merge

def get_config_schema():
    with open(os.path.dirname(__file__) +  '/config_schema.json') as f:
        return json.loads(f.read())

skema_defaults = '''
ObjectId: Any
DateTime: Any
Date: Any
Time: Any
Json: Any
'''

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


def get_skema(config, here='./'):
    if "skema_path" in config:
        path = here + config["skema_path"]
        with open(path) as f:
            return f.read() + skema_defaults
    if config.get("skema"):
        return config.get("skema") + skema_defaults
    if config.get("skema_url"):
        r = requests.get(config.get("skema_url"), stream=True)
        skema = ""
        while 1:
            buf = r.raw.read(16 * 1024)
            if not buf:
                break
            skema += buf.decode()
        return skema + skema_defaults

