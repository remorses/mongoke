


from mongoke.templates.support import join_yields, repr_eval_dict



def test_join_yields():
    @join_yields(', ')
    def op():
        yield '1'
        yield '2'
        yield '3'
    x = op()
    print(x)

def test_repr_eval_dict():
    x = {
        'eq': {
            'not': "parent['x'] in headers[\"ciao\"]"
        }
    }
    print(repr_eval_dict(x, '    ').lstrip())