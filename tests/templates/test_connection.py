# import pytest
import mongodb_streams
# import mock
# from generated.generated.resolvers.support import connection_resolver


# @pytest.mark.asyncio
# async def test_1():
#     with mock.patch('mongodb_streams.find', ) as m:
#         result = await connection_resolver(
#             collection=None,
#             where={},
#             cursorField='timestamp',  # needs to exist always at least one, the fisrst is the cursorField
#             pagination={'last': 3},
#             scalar_name='String',
#             pipeline=[]
#         )