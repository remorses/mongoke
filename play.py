import yaml
import json
s = '''
ciao:
    expr: > 
        parent['xxx'] == 'ciao'
        and parent is not None
'''

# print(yaml.load(s)['ciao']['expr'])
x = yaml.load(open('pr_conf.yaml'))
y = json.dumps(x, indent=4)
print(y)skema