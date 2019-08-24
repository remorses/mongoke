
import jwt
from bson import ObjectId
from aiohttp import web
from ..logger import logger

JWT_ALGORITHM = 'HS256'

async def jwt_middleware(app, handler):
    async def middleware(request: web.Request):
        request.jwt_payload = {}
        jwt_token = request.headers.get('Authorization', None)
        if jwt_token:
            try:
                payload = jwt.decode(jwt_token, verify=False, algorithms=[JWT_ALGORITHM])
            except (jwt.DecodeError, jwt.ExpiredSignatureError) as e:
                # return to_response(*fail('Token is invalid'))
                logger.error(e)
                return await handler(request)
            else:
                request.jwt_payload = payload
                # db = request.app.db
                # user = await db.businessUsers.find_one({ # TODO not all the user details
                #     '_id': ObjectId(payload['user_id'])
                # })
                # if not user:
                #     logger.error(f'error no user in valid jwt for {payload["user_id"]}')
                #     return to_response(*fail('User profile is compromised, valid jwt user not in db'))
                # else:
                #     request.user = {**user, 'user_id': user['_id']}
        else:
            logger.debug('no Authorization header')
        return await handler(request)
    return middleware

