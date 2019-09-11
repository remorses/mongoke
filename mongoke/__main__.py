from .generate_from_config import generate_from_config
import sys
import yaml

if __name__ == "__main__":
    arg = sys.argv[-1]
    config = yaml.safe_load(open(arg).read())
    generate_from_config(config)
