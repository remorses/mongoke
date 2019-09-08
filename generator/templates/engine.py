engine = '''
from tartiflette import Engine
from typing import *
from .generated.logger import logger

async def my_error_coercer(
    exception: Exception, error: Dict[str, Any]
) -> Dict[str, Any]:
    # error["extensions"]["type"] = "custom_exception"
    logger.exception(exception)
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