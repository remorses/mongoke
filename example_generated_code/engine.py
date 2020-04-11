
import os
import sys
from tartiflette import Engine
from typing import *
from .generated.logger import logger
from tartiflette.types.exceptions.tartiflette import TartifletteError
import traceback

def read(path):
    if path:
        with open(path) as f:
            return f.read()
    return ''

async def my_error_coercer(
    exception: TartifletteError, error: Dict[str, Any]
) -> Dict[str, Any]:

    ex = exception.original_error
    if ex and os.getenv('DEBUG'):
        trace = "\n".join(
            traceback.format_exception(etype=type(ex), value=ex, tb=ex.__traceback__)
        )
        logger.error(trace)
    else:
        logger.error(exception)
    message = str(ex) if ex else ''
    if message:
        return {**error, 'message': f'Error in Mongoke server: {message}'}
    else:
        return error


class CustomEngine(Engine):
    async def cook(
        self,
        sdl: Union[str, List[str]] = None,
        error_coercer=my_error_coercer,
        custom_default_resolver: Optional[Callable] = None,
        custom_default_type_resolver: Optional[Callable] = None,
        modules: Optional[Union[str, List[str], List[Dict[str, Any]]]] = None,
        schema_name: str = None,
    ):
        try:
            await super().cook(
                sdl, my_error_coercer, custom_default_resolver, custom_default_type_resolver, modules, schema_name
            )
        except Exception as e:
            print('ERROR parsing graphql schema:')
            print(e)
            raise e

