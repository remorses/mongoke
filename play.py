import yaml

s = '''
ciao:
    expr: > 
        parent['xxx'] == 'ciao'
        and parent is not None
'''

# print(yaml.load(s)['ciao']['expr'])
x = yaml.load(open('spec_conf.yaml'))
y = yaml.dump(x)
print(y)