import os
import coloredlogs
import logging

logger = logging.getLogger(__name__)

LEVEL = 'DEBUG' if os.getenv('DEBUG') else 'INFO'

coloredlogs.DEFAULT_FIELD_STYLES = {'asctime': {'color': 'white'}, 'hostname': {'color': 'white'}, 'levelname': {
    'color': 'white', 'bold': True}, 'name': {'color': 'white'}, 'programname': {'color': 'white'}}
coloredlogs.install(
    fmt='%(asctime)s %(module)s %(levelname)s %(message)s', level=LEVEL, logger=logger)
