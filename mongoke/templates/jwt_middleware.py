# jwt_header
# jwt_sheme
# jwt_required
# jwt_secret
# jwt_algorithms   



jwt_middleware = '''
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.responses import Response
import jwt
from bson import ObjectId
from .generated.logger import logger

JWT_ALGORITHMS = ${{ repr(jwt_algorithms) }}


class JwtMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, handler):
        request.state.jwt_payload = {}
        jwt_token = (
            request.headers.get("${{ jwt_header }}", "").replace("${{ jwt_sheme }}", "").strip()
        )
        if not jwt_token:
            ${{
                indent_to('            ', """
                return Response(status_code=401, content='Missing authorization token')
                """) if jwt_required else indent_to('            ', 'return await handler(request)')
            }}
        try:
            payload = jwt.decode(
                jwt_token, verify=${{ repr(bool(jwt_required)) }}, secret=${{ repr(jwt_secret) }}, algorithms=[JWT_ALGORITHMS]
            )
        except (jwt.InvalidTokenError) as exc:
            logger.error("Cannot decode authorization token, " + str(exc))
            if jwt_required:
                logger.error("returning error 403 as jwt is required")
                msg = "Invalid authorization token, " + str(exc)
                return Response(status_code=403, content=msg)
        except Exception as exc:
            logger.error("Cannot decode authorization token, " + str(exc))
        else:
            request.state.jwt_payload = payload
        return await handler(request)

'''



_jwt_middleware = '''
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
                """) if jwt_required else indent_to('            ', 'return await handler(request)')
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