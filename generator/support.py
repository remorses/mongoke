import json
import os.path
from funcy import pluck, count

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