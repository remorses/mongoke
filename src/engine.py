
from tartiflette import Engine
from typing import *
from .generated.logger import logger
import traceback
from tartiflette.types.exceptions.tartiflette import TartifletteError
async def my_error_coercer(
    ex: Exception, error: Dict[str, Any]
) -> Dict[str, Any]:
    # error["extensions"]["type"] = "custom_exception"
    logger.warn(dir(ex))
    ex = ex.original_error
    trace = '\n'.join(traceback.format_exception(etype=type(ex), value=ex, tb=ex.__traceback__))
    logger.error(trace)
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
