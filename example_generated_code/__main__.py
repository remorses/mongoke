
from .make_app import make_app
try:
    app = make_app()
except Exception as e:
    print(e)
