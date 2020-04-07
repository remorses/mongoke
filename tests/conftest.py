import asyncio
import mongodb_streams
import pytest
from bson import ObjectId
from pymongo.mongo_client import MongoClient
from pymongo.database import Database
from pymongo.collection import Collection
from motor.motor_asyncio import AsyncIOMotorClient
from starlette.testclient import TestClient
from example_generated_code.make_app import make_app
from unittest.mock import _Call, call
from asynctest import mock
from prtty import pretty
from starlette.routing import Lifespan, Router


DB_URL = "mongodb://localhost/testdb"


@pytest.fixture
def mongo(event_loop):
    mongo = AsyncIOMotorClient(DB_URL)
    return mongo


@pytest.fixture
async def db(mongo):
    mongo: MongoClient = mongo.delegate
    return mongo.get_database("testingdatabase")


@pytest.fixture
async def asyncdb(event_loop, mongo):
    db = mongo.get_database("testingdatabase")
    yield db
    # for c in await db.list_collection_names():
    #     await db.drop_collection(c)


@pytest.fixture
def client(asyncdb):
    app = make_app(asyncdb)

    with TestClient(app) as client:
        yield client


@pytest.fixture
def query(client, db):
    def func(query, variables={}):
        r = client.post("/", json=dict(query=query, variables=variables))
        return r.json()

    # await client.wait_startup()
    # client.__enter__()
    return func
    # client.__exit__(None)


@pytest.fixture
def users(db):
    users: Collection = db.users
    yield users
    users.delete_many({})
