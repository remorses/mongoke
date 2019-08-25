from funcy import lfilter, post_processing, pluck, count

def zip_pluck(d, keys, enumerate=False):
    args = [pluck(k, d) for k in keys]
    if enumerate:
        args = [count(), *args]
    return zip(*args)

join_yields = lambda separator='': post_processing(lambda parts: separator.join(list(map(str, parts))))

def remove_indentation(string: str):
    lines = lfilter(bool, string.split('\n'))
    print(lines)
    base_indent = min([len(line) - len(line.lstrip()) for line in lines])
    lines = [x[base_indent:] for x in lines]
    return '\n'.join(lines)

def indent_to(indentation, string):
    string = remove_indentation(string)
    return '\n'.join([indentation + line for line in string.split('\n')])


if __name__ == '__main__':
    join_yields()
    print('hello')
    x = remove_indentation("""
    ciao x
        come va
    """)
    print(x)
    x = '''
    xxx
        yy
    '''
    print(indent_to('....', x))