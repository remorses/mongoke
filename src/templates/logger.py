
logger = '''
import logging
import coloredlogs

logging_format = "[%(asctime)s] [%(levelname)s] "
logging_format += "[%(module)s] "
logging_format += "%(message)s"

logger = logging.getLogger(__name__)
if not len(logger.handlers):
    ch = logging.StreamHandler()
    ch.setFormatter(logging.Formatter(logging_format))
    ch.setLevel(logging.DEBUG)
    logger.addHandler(ch)
    logger.setLevel(logging.DEBUG)

coloredlogs.DEFAULT_FIELD_STYLES = {'asctime': {'color': 'white'}, 'hostname': {'color': 'white'}, 'levelname': {'color': 'white', 'bold': True}, 'module': {'color': 'white', 'bold': True}, 'name': {'color': 'white'}, 'programname': {'color': 'white'}}
coloredlogs.install(fmt=logging_format, level=logging.DEBUG, logger=logger)
'''