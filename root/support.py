from funcy import pluck

def zip_pluck(d, *keys):
    return zip(*[pluck(k, d) for k in keys])