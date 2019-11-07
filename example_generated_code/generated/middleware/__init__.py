
import jwt
from bson import ObjectId
from aiohttp import web
from ..logger import logger

JWT_ALGORITHMS = ['H256']


async def jwt_middleware(app, handler):
    async def middleware(request: web.Request):
        request.jwt_payload = {}
        jwt_token = request.headers.get('Authorization', '').replace('Bearer', '').strip()
        if not jwt_token:
            return await handler(request)
        try:
            payload = jwt.decode(jwt_token, verify=False, secret=None, algorithms=[JWT_ALGORITHMS])
        except (jwt.InvalidTokenError) as exc:
            logger.exception(exc, exc_info=exc)
            msg = 'Invalid authorization token, ' + str(exc)
            raise web.HTTPForbidden(reason=msg)
        else:
            request.jwt_payload = payload
        return await handler(request)
    return middleware
