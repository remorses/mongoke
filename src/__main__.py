
import os
import urllib.parse
from aiohttp import web
from motor.motor_asyncio import AsyncIOMotorClient
from tartiflette_aiohttp import register_graphql_handlers
import asyncio

import src.generated.resolvers
import src.generated.scalars
from src.generated.middleware import jwt_middleware

here = os.path.dirname(os.path.abspath(__file__))

DB_URL = "" or None

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
    app.on_startup.append(lambda app: context.update({'loop': asyncio.get_event_loop()}))
    web.run_app(app)

run()
