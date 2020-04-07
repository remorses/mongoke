engine = '''
import os
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
    if ex and not os.getenv('HIDE_ERRORS_TRACEBACK'):
        trace = "\\n".join(
            traceback.format_exception(etype=type(ex), value=ex, tb=ex.__traceback__)
        )
        logger.error(trace)
    else:
        logger.error(exception)
    better_error = {**error, 'message': f'Error in Mongoke server: {str(ex)}'}
    return better_error


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
        await super().cook(
            sdl, my_error_coercer, custom_default_resolver, custom_default_type_resolver, modules, schema_name
        )

'''
