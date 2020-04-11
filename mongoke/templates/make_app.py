make_app = '''

import os
import sys
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

${{'\\n'.join([f'import {root_dir_name}.generated.resolvers.{name}' for name in resolver_names])}}
import ${{root_dir_name}}.generated.scalars



GRAPHIQL_DEFAULT_QUERY = """
# welcome to mongoke
# ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭╮
# ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱┃┃
# ╭╮╭┳━━┳━╮╭━━┳━━┫┃╭┳━━╮
# ┃╰╯┃╭╮┃╭╮┫╭╮┃╭╮┃╰╯┫┃━┫
# ┃┃┃┃╰╯┃┃┃┃╰╯┃╰╯┃╭╮┫┃━┫
# ╰┻┻┻━━┻╯╰┻━╮┣━━┻╯╰┻━━╯
# ╱╱╱╱╱╱╱╱╱╭━╯┃
# ╱╱╱╱╱╱╱╱╱╰━━╯

# Mongoke generates graphql queries for your mongodb data
# To get started try a one document query like this
# query OneDocument {
#     YourTypeName {
#         yourField
#     }
# }


# you can use the where argument to filter documents
# query OneDocument {
#     YourTypeName(where: { color: { eq: "red" } }) {
#         yourField
#     }
# }


# try now querying multiple documents
# the objects are returned with the pageInfo data, useful for pagination
# the sorting direction and pageInfo.endCursor are based on the 'cursorField' argument, this can be any scalar field of your document
# query MultipleDocuments  {
#     YourTypeNames(cursorField: color) {
#         nodes {
#             color
#         }
#         pageInfo {
#             hasNextPage
#             endCursor # use this as the after field for the next query
#         }
#     }
# }

# you can use the pageInfo.endCursor to fetch the next documents
# query MultipleDocuments  {
#     YourTypeNames(first: 10, after: "red", cursorField: color) {
#         nodes {
#             field
#         }
#         pageInfo {
#             hasNextPage
#             endCursor # use this as the after field for the next query
#         }
#     }
# }

"""


MONGOKE_BASE_PATH = os.getenv("MONGOKE_BASE_PATH", "/")
DISABLE_GRAPHIQL = bool(os.getenv("DISABLE_GRAPHIQL", False))
GRAPHIQL_DEFAULT_JWT = os.getenv("GRAPHIQL_DEFAULT_JWT", "")
GRAPHIQL_QUERY = os.getenv("GRAPHIQL_DEFAULT_QUERY", "") or read(
    os.getenv("GRAPHIQL_DEFAULT_QUERY_FILE_PATH", "")
) or GRAPHIQL_DEFAULT_QUERY
DB_URL = os.getenv("DB_URL") or "${{db_url}}" or None


here = os.path.dirname(os.path.abspath(__file__))
sdl_dir = f"{here}/generated/sdl/"
sdl_files = sorted(os.listdir(sdl_dir))
#  print(sdl_files)
sdl_files = [sdl_dir + f for f in sdl_files]


class CatchAll(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, handler):
        request.scope["path"] = MONGOKE_BASE_PATH # TODO subscriptions path
        return await handler(request)

engine = CustomEngine(
    sdl=sdl_files, modules=[ApolloFederationPlugin(engine_sdl=sdl_files)]
)

def make_app(db: AsyncIOMotorClient=None):
    try:
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
        app = CORSMiddleware(app, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"], )
        app = JwtMiddleware(app,)
        # app = CatchAll(app,)
        app = ServerErrorMiddleware(app,)
        return app
    except Exception as e:
        print('got an error starting the Mongoke server:')
        print(e)
        return







'''
