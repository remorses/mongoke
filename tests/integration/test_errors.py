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

LETTERS = "abcdefghilmnopqrstuvz"
_ = mock.ANY


def test_no_error(query, users: Collection):
    assert not list(users.find({}))
    q = """
        {
            User {
                _id
                name
            }
        }

      """
    res = query(q)
    pretty(res)
    assert not res.get("errors")


def test_objectid_error(query, users: Collection):
    assert not list(users.find({}))
    q = """
        {
            User(where: {_id: {eq: "kjh"}}) {
                _id
                name
            }
        }

      """
    res = query(q)
    pretty(res)
    assert res["errors"]
    pretty(res["errors"][0]["message"])
