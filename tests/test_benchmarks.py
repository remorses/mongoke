
from src.__main__ import build

async def make_query(client, query):
    

def test_1(aiohttp_client, benchmark):
    app = build()
    client = await aiohttp_client(app)