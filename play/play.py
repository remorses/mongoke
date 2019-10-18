import yaml
import json
import jsonschema
import os.path

print(os.path.dirname(__file__))

s = '''
ciao:
    expr: > 
        parent['xxx'] == 'ciao'
        and parent is not None
'''

# print(yaml.load(s)['ciao']['expr'])
x = yaml.load(open('../pr_conf.yaml'))
jsonschema.validate(x, json.load(open('configuration_schema.json')))
# y = json.dumps(x, indent=4)
# print(y)