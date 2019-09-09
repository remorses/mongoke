import asyncio
import mongodb_streams
import pytest
from generated.__main__ import build
from unittest.mock import _Call, call
from asynctest import mock
from pprint import pprint

from aiohttp.test_utils import TestClient, TestServer

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
        pprint(r)
        pprint(m.call_args_list)
        # m.assert_called_with(_, where={'username': {'$eq': 'ciao'}}, pipeline=_)


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
        pprint(r, indent=4)
        print(m.call_args)
        
        # m.assert_called_with(_, where={'username': {'$eq': 'ciao'}}, pipeline=_)
