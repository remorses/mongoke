import asyncio
import mongodb_streams
import pytest
from aiohttp.test_utils import TestClient, TestServer
from generated.__main__ import build
from unittest.mock import _Call, call
from asynctest import mock
from prtty import pretty

_ = mock.ANY


@pytest.fixture()
async def query():
    print('quering')
    loop = asyncio.get_event_loop()
    loop.set_debug(True)
    app = build(db=mock.MagicMock(name='db'),)
    async with TestClient(TestServer(app), loop=loop) as client:
        async def func(query, variables={}):
            r = await client.post('/', json=dict(query=query, variables=variables))
            return await r.json()
        yield func


@pytest.mark.asyncio
async def test_single_resolver(query):
    with mock.patch('mongodb_streams.find_one', ) as m:
        m.side_effect = [
            dict(username='hello'),
            dict(name='name')
        ]
        r = await query('''
        {
            bot(where: {username: {eq: "ciao"}}) {
                username
                user {
                    name
                }
            }
        }
        ''')
        pretty(r)
        pretty(m.call_args_list)
        assert r == {
            "data": {
                "bot": {
                    "username": "hello",
                    "user": {
                        "name": "name"
                    }
                }
            }
        }

@pytest.mark.asyncio
async def test_many_relation(query):
    with mock.patch('mongodb_streams.find_one', ) as m:
        m.side_effect = [
            [dict(username='hello'),],
            [dict(value=89, timestamp=34)]
        ]
        r = await query('''
            {
                bots(
                    last: 50,
                ) {
                    nodes {
                        username
                        _id
                        likes_over_time(
                            first: 20
                            cursorField: value
                        ) {
                            nodes {
                                value
                                timestamp
                            }
                        }
                    }
                }
            }
        ''')
        pretty(r)
        pretty(m.call_args_list)
        


@pytest.mark.asyncio
async def test_many_resolver(query):
    with mock.patch('mongodb_streams.find', ) as m:
        bots = [dict(_id=str(i), username=str(i)) for i in range(20)]
        m.return_value = bots
        r = await query('''
        {
            bots(first: 3) {
                nodes {
                    username
                }
            }
        }
        ''')
        pretty(r, )
        print(m.call_args)
        
        # m.assert_called_with(_, where={'username': {'$eq': 'ciao'}}, pipeline=_)
