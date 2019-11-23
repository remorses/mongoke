
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.responses import Response
import jwt
from bson import ObjectId
from .generated.logger import logger

JWT_ALGORITHMS = ['H256']


class JwtMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request, handler):
        request.state.jwt_payload = {}
        jwt_token = (
            request.headers.get("Authorization", "").replace("Bearer", "").strip()
        )
        if not jwt_token:
            return await handler(request)
        try:
            payload = jwt.decode(
                jwt_token, verify=False, secret=None, algorithms=[JWT_ALGORITHMS]
            )
        except (jwt.InvalidTokenError) as exc:
            logger.exception(exc, exc_info=exc)
            msg = "Invalid authorization token, " + str(exc)
            return Response(status_code=403, content=msg)
        else:
            request.state.jwt_payload = payload
        return await handler(request)

