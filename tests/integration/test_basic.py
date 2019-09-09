import asyncio
import mongodb_streams
import pytest
from generated.__main__ import build
from unittest.mock import _Call, call
from  asynctest import mock

from aiohttp.test_utils import TestClient, TestServer, loop_context


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
async def test_1(query):
    with mock.patch('mongodb_streams.find_one', ) as m:
        m.return_value = dict(username='hello')
        r = await query('''
        {
            bot(where: {username: {eq: "ciao"}}) {
                username
            }
        }
        ''')
        print(m.call_args)
        assert m.call_args == (mock.ANY, {'username': {'$eq': 'ciao'}}, mock.ANY)
        # m.assert_called_with(pipeline=mock.ANY)
