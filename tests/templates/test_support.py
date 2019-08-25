


from generator.templates.support import join_yields



def test_join_yields():
    @join_yields(', ')
    def op():
        yield 1
        yield '2'
        yield '3'
    x = op()
    print(x)