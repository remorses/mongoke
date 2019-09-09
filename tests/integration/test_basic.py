import asyncio
import mongodb_streams
import pytest
from generated.__main__ import build
from unittest.mock import _Call
import asynctest
import mock
from aiohttp.test_utils import TestClient, TestServer, loop_context


@pytest.fixture()
async def query():
    print('quering')
    loop = asyncio.get_event_loop()
    loop.set_debug(True)
    app = build(db=mock.MagicMock(),)
    async with TestClient(TestServer(app), loop=loop) as client:
        async def func(query, variables={}):
            r = await client.post('/', json=dict(query=query, variables=variables))
            return await r.json()
        yield func
    
@pytest.fixture
def afind_one():
    print('mocking')
    m = asynctest.mock.patch('mongodb_streams.find_one', ).__enter__()
    m.return_value = dict()
    yield m

@pytest.mark.asyncio
async def test_1(afind_one, query):
    afind_one.return_value = dict()
    r = await query('''
    {
        bot {
            username
        }
    }
    ''')
    print(r)
