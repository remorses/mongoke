main = '''
import os
import aiohttp_cors
import urllib.parse
from aiohttp import web
from motor.motor_asyncio import AsyncIOMotorClient
from tartiflette_aiohttp import register_graphql_handlers
import asyncio

${{'\\n'.join([f'import {root_dir_name}.generated.resolvers.{name}' for name in resolver_names])}}
import ${{root_dir_name}}.generated.scalars
from ${{root_dir_name}}.generated.middleware import jwt_middleware

here = os.path.dirname(os.path.abspath(__file__))

DB_URL = "${{db_url}}" or None

def run():
    app = web.Application(middlewares=[jwt_middleware])
    db = AsyncIOMotorClient(DB_URL)
    app.db = db
    context = {
        'db': db,
        'app': app,
        'loop': None,
    }
    app = register_graphql_handlers(
        app=app,
        engine_sdl=f'{here}/generated/sdl/',
        executor_context=context,
        executor_http_endpoint='/',
        executor_http_methods=['POST', 'GET',],
        graphiql_enabled=True
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
    web.run_app(app)

run()
'''