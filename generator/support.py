import json
import os.path
from funcy import pluck

def touch(filename, data):
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, "w") as f:
        f.write(data)

pretty = lambda x: print(json.dumps(x, indent=4, default=str))

def zip_pluck(d, keys):
    return zip(*[pluck(k, d) for k in keys])