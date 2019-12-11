import asyncio
import mongodb_streams
import pytest
from starlette.testclient import TestClient
from example_generated_code.make_app import make_app
from unittest.mock import _Call, call
from asynctest import mock
from prtty import pretty
from starlette.routing import Lifespan, Router

LETTERS = "abcdefghilmnopqrs"
_ = mock.ANY


@pytest.fixture
def client():
    app = make_app(db=mock.MagicMock(name="db"))

    with TestClient(app) as client:
        yield client
    del app


@pytest.fixture
def query(client):
    def func(query, variables={}):
        r = client.post("/", json=dict(query=query, variables=variables))
        return r.json()

    # await client.wait_startup()
    # client.__enter__()
    return func
    # client.__exit__(None)


def test_get_user(query):
    q = """
        {
            User {
                _id
                name
            }
        }

      """
    res = query(q)
    print(res)


def test_single_resolver(query):
    with mock.patch("mongodb_streams.find_one") as m:
        m.side_effect = [dict(surname="xxx")]
        r = query(
            """
            {
                User(where: {surname: {eq: "xxx"}}) {
                    _id
                    name
                    surname
                }
            }
        """
        )
        pretty(r)
        pretty(m.call_args_list)
        assert r["data"]["User"]["surname"] == "xxx"


def test_many_users(query):
    with mock.patch("mongodb_streams.find") as m:
        m.side_effect = [[dict(surname="xxx")], [dict(surname="xxx")]]
        r = query(
            """
                {
                    Users(where: { surname: { eq: "xxx" } }) {
                        nodes {
                            surname
                        }
                    }
                }

            """
        )
        pretty(r)
        pretty(m.call_args_list)
        nodes = r["data"]["Users"]["nodes"]
        assert len(nodes) == 1
        assert nodes[0]["surname"] == "xxx"


def test_many_resolver(query):
    with mock.patch("mongodb_streams.find") as m:
        bots = [dict(_id=str(i), username=str(i)) for i in range(20)]
        m.return_value = bots
        r = query(
            """
            {
                Users(where: { surname: { eq: "xxx" } }) {
                    nodes {
                        surname
                        friends {
                            nodes {
                                surname
                            }
                        }
                    }
                }
            }
            """
        )
        pretty(r)
        print(m.call_args)


def test_cursor_field(query):
    with mock.patch("mongodb_streams.find") as m:
        xs = [dict(surname=i) for i in iter(LETTERS)]
        m.return_value = xs
        r = query(
            """
            {
                Users(cursorField: surname) {
                    nodes {
                        surname
                    }
                }
            }
            """
        )
        pretty(r)
        print(m.call_args)
        nodes = r["data"]["Users"]["nodes"]
        assert len(nodes) == len(LETTERS)
        assert sorted(nodes, key=lambda x: x["surname"]) == nodes

        # m.assert_called_with(_, where={'username': {'$eq': 'ciao'}}, pipeline=_)
