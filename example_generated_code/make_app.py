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

import example_generated_code.generated.resolvers.User
import example_generated_code.generated.resolvers.Users
import example_generated_code.generated.resolvers.Human
import example_generated_code.generated.resolvers.Humans
import example_generated_code.generated.resolvers.Task_events
import example_generated_code.generated.resolvers.User_friends
import example_generated_code.generated.resolvers.User_likes_over_time
import example_generated_code.generated.resolvers.User_father
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


class CatchAll(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, handler):
        request.scope["path"] = MONGOKE_BASE_PATH  # TODO subscriptions path
        return await handler(request)


engine = CustomEngine(
    sdl=sdl_files, modules=[ApolloFederationPlugin(engine_sdl=sdl_files)]
)


def make_app(db: AsyncIOMotorClient = None):
    if not db:
        db = AsyncIOMotorClient(DB_URL).get_database()
    graphiql = GraphiQL(
        # path=MONGOKE_BASE_PATH,
        default_headers={"Authorization": "Bearer " + GRAPHIQL_DEFAULT_JWT}
        if GRAPHIQL_DEFAULT_JWT
        else {},
        default_query=GRAPHIQL_QUERY,
    )

    context = {"db": db, "loop": None}

    app = TartifletteApp(
        context=context,
        engine=engine,
        path=MONGOKE_BASE_PATH,
        graphiql=graphiql if not DISABLE_GRAPHIQL else False,
    )
    app = CORSMiddleware(
        app, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"]
    )
    app = JwtMiddleware(app)
    # app = CatchAll(app,)
    app = ServerErrorMiddleware(app)
    return app

