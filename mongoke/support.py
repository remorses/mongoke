import json
import os.path
import requests
from funcy import pluck, count, merge

def touch(filename, data):
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, "w") as f:
        f.write(data)

pretty = lambda x: print(json.dumps(x, indent=4, default=str))

def zip_pluck(d, keys, enumerate=False):
    args = [pluck(k, d) for k in keys]
    if enumerate:
        args = [count(), *args]
    return zip(*args)

def get_skema(config):
    if 'skema_path' in config:
        with open(config['skema_path']) as f:
            return f.read()
    if config.get('skema'):
        return config.get('skema')
    if config.get('skema_url'):
        r = requests.get(config.get('skema_url'), stream=True)
        skema = ''
        while 1:
            buf = r.raw.read(16*1024)
            if not buf:
                break
            skema += buf.decode()
        return skema

