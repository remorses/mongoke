

import os
from tartiflette import Resolver, Engine
from tartiflette_asgi import TartifletteApp, GraphiQL
from tartiflette_plugin_apollo_federation import ApolloFederationPlugin
from starlette.middleware.cors import CORSMiddleware
from starlette.requests import Request
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.middleware.errors import ServerErrorMiddleware
from motor.motor_asyncio import AsyncIOMotorClient
from .middleware import JwtMiddleware
from .engine import CustomEngine, read

import example_generated_code.generated.resolvers.user
import example_generated_code.generated.resolvers.users
import example_generated_code.generated.resolvers.human
import example_generated_code.generated.resolvers.humans
import example_generated_code.generated.resolvers.task_events
import example_generated_code.generated.resolvers.user_friends
import example_generated_code.generated.resolvers.user_likes_over_time
import example_generated_code.generated.resolvers.user_father
import example_generated_code.generated.scalars

MONGOKE_BASE_PATH = os.getenv("MONGOKE_BASE_PATH", "/")
DISABLE_GRAPHIQL = bool(os.getenv("DISABLE_GRAPHIQL", False))
GRAPHIQL_DEFAULT_JWT = os.getenv("GRAPHIQL_DEFAULT_JWT", "")
GRAPHIQL_QUERY = os.getenv("GRAPHIQL_DEFAULT_QUERY", "") or read(
    os.getenv("GRAPHIQL_DEFAULT_QUERY_FILE_PATH", "")
)
DB_URL = os.getenv("DB_URL") or "" or None


here = os.path.dirname(os.path.abspath(__file__))
sdl_dir = f"{here}/generated/sdl/"
sdl_files = sorted(os.listdir(sdl_dir))
#  print(sdl_files)
sdl_files = [sdl_dir + f for f in sdl_files]


def make_app():
    graphiql = GraphiQL(
        # path=MONGOKE_BASE_PATH,
        default_headers={"Authorization": "Bearer " + GRAPHIQL_DEFAULT_JWT}
        if GRAPHIQL_DEFAULT_JWT
        else {},
        default_query=GRAPHIQL_QUERY,
    )

    engine = CustomEngine(
        sdl=sdl_files, modules=[ApolloFederationPlugin(engine_sdl=sdl_files)]
    )

    db: AsyncIOMotorClient = AsyncIOMotorClient(DB_URL).get_database()

    context = {"db": db, "loop": None}

    app = TartifletteApp(
        context=context,
        engine=engine,
        path=MONGOKE_BASE_PATH,
        graphiql=graphiql if not DISABLE_GRAPHIQL else False,
    )
    return app

class CatchAll(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, handler):
        request.scope["path"] = MONGOKE_BASE_PATH
        return await handler(request)

app = make_app()
app = CORSMiddleware(app, allow_origins=["*"], allow_methods=["*"])
app = JwtMiddleware(app,)
app = ServerErrorMiddleware(app,)
app = CatchAll(app,)


