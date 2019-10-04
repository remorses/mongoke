import os.path
import fire
from .generate_from_config import generate_from_config
from .checksum import make_checksum, existent_checksum
import sys
import yaml

def main(path, force=False):
    config = yaml.safe_load(open(path).read())
    config_path = os.path.abspath(os.path.dirname(path))
    if not force:
        checksum = make_checksum(config, config_path + '/')
        if existent_checksum(config, ) == checksum:
            print('already generated')
            return
    generate_from_config(config, config_path + '/')