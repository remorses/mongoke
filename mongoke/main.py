import os.path
import fire
from .generate_from_config import generate_from_config
import sys
import yaml

def main(path):
    config = yaml.safe_load(open(path).read())
    config_path = os.path.abspath(os.path.dirname(path))
    generate_from_config(config, config_path + '/')