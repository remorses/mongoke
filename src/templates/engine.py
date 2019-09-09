engine = '''
from tartiflette import Engine
from typing import *
from .generated.logger import logger
from tartiflette.types.exceptions.tartiflette import TartifletteError
import traceback

async def my_error_coercer(
    exception: TartifletteError, error: Dict[str, Any]
) -> Dict[str, Any]:

    ex = exception.original_error
    if ex:
        trace = '\\n'.join(traceback.format_exception(etype=type(ex), value=ex, tb=ex.__traceback__))
        logger.error(trace)
    else:
        logger.error(exception)
    return error

class CustomEngine(Engine):
    async def cook(
        self,
        sdl: Union[str, List[str]],
        error_coercer: Callable[[Exception], dict] = None,
        custom_default_resolver: Optional[Callable] = None,
        modules: Optional[Union[str, List[str]]] = None,
        schema_name: str = "default",
    ):
        await super().cook(
            sdl,
            my_error_coercer,
            custom_default_resolver,
            modules,
            schema_name
        )
'''