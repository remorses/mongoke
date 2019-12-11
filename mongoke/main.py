import os.path
import fire
from .generate_from_config import generate_from_config
from .checksum import make_checksum, existent_checksum
from .support import get_config_schema
import jsonschema
import sys
import yaml
from mongoke.support import download_file



def main(path, force=False, generated_path=None):
    if os.getenv('MONGOKE_CONFIG_URL'):
        config = download_file(os.getenv('MONGOKE_CONFIG_URL'))
    else:
        config = open(path).read()
    config = yaml.safe_load(config)
    jsonschema.validate(config, get_config_schema())
    config_path = os.path.abspath(os.path.dirname(path))
    if not force:
        checksum = make_checksum(config, config_path + '/')
        if existent_checksum(config, ) == checksum:
            print('already generated')
            return
    generate_from_config(config, config_path + '/', generated_path)