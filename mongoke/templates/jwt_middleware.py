# jwt_header
# jwt_sheme
# jwt_required
# jwt_secret
# jwt_algorithms   
jwt_middleware = '''
import jwt
from bson import ObjectId
from aiohttp import web
from ..logger import logger

JWT_ALGORITHMS = ${{ repr(jwt_algorithms) }}


async def jwt_middleware(app, handler):
    async def middleware(request: web.Request):
        request.jwt_payload = {}
        jwt_token = request.headers.get(${{ repr(jwt_header) }}, '').replace(${{ repr(jwt_sheme) }}, '').strip()
        if not jwt_token:
            ${{
                indent_to('            ', """
                raise web.HTTPUnauthorized(
                    reason='Missing authorization token',
                )
                """) if jwt_required else indent_to('            ', 'pass')
            }}
        try:
            payload = jwt.decode(jwt_token, verify=${{ repr(bool(jwt_required)) }}, secret=${{ repr(jwt_secret) }}, algorithms=[JWT_ALGORITHMS])
        except (jwt.InvalidTokenError) as exc:
            logger.exception(exc, exc_info=exc)
            msg = 'Invalid authorization token, ' + str(exc)
            raise web.HTTPForbidden(reason=msg)
        else:
            request.jwt_payload = payload
        return await handler(request)
    return middleware
'''