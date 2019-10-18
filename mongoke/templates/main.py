main = '''
import os
import aiohttp_cors
import urllib.parse
from aiohttp import web
from motor.motor_asyncio import AsyncIOMotorClient
from tartiflette_aiohttp import register_graphql_handlers
from tartiflette_plugin_apollo_federation import ApolloFederationPlugin
import asyncio

from .engine import CustomEngine, read
${{'\\n'.join([f'import {root_dir_name}.generated.resolvers.{name}' for name in resolver_names])}}
import ${{root_dir_name}}.generated.scalars
from ${{root_dir_name}}.generated.middleware import jwt_middleware

DB_URL = os.getenv('DB_URL') or "${{db_url}}" or None
PORT = 80

here = os.path.dirname(os.path.abspath(__file__))
sdl_dir = f'{here}/generated/sdl/'
sdl_files = sorted(os.listdir(sdl_dir))
# print(sdl_files)
sdl_files = [sdl_dir + f for f in sdl_files]

GRAPHIQL_QUERY = os.getenv("GRAPHIQL_DEFAULT_QUERY", "") or read(os.getenv('GRAPHIQL_DEFAULT_QUERY_FILE_PATH', ''))

def build(db):
    app = web.Application(middlewares=[jwt_middleware])
    app.db = db
    context = {
        'db': db,
        'app': app,
        'loop': None,
    }
    app = register_graphql_handlers(
        app=app,
        engine=CustomEngine(),
        engine_sdl=sdl_files,
        executor_context=context,
        executor_http_endpoint='/',
        executor_http_methods=['POST', 'GET',],
        engine_modules=[
            ApolloFederationPlugin(engine_sdl=sdl_files)
        ],
        graphiql_enabled=os.getenv("DISABLE_GRAPHIQL", True),
        graphiql_options={
            "default_query": GRAPHIQL_QUERY,
            "default_variables": {},
            "default_headers": {
                "Authorization": "Bearer " + os.getenv("GRAPHIQL_DEFAULT_JWT", "")
            }
            if os.getenv("GRAPHIQL_DEFAULT_JWT")
            else {},
        },
    )
    cors = aiohttp_cors.setup(app, defaults={
        "*": aiohttp_cors.ResourceOptions(
                allow_credentials=True,
                expose_headers="*",
                allow_headers="*",
            )
    })
    for route in list(app.router.routes()):
        cors.add(route)
    async def on_startup(app):
        context.update({'loop': asyncio.get_event_loop()})
    app.on_startup.append(on_startup)
    return app

if __name__ == '__main__':
    db: AsyncIOMotorClient = AsyncIOMotorClient(DB_URL).get_database()
    web.run_app(build(db), port=PORT)


'''