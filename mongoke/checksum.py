import json
import os.path
from .support import get_types_schema
import hashlib


def existent_checksum(config, ):
    root_dir_path = config.get("root_dir_path", "generated")
    checksum_path = os.path.join(root_dir_path, 'checksum')
    if os.path.exists(checksum_path):
        checksum = open(checksum_path).read().strip()
        return checksum
    else:
        return None

def make_checksum(config, config_path):
    skema = get_types_schema(config, here=config_path)
    config = {**config, 'skema': skema}
    config = json.dumps(config, sort_keys=True)
    return hashlib.md5(config.encode()).hexdigest()