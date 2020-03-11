#!/bin/sh
# pwd
python -m mongoke $1 && uvicorn generated.__main__:app --host 0.0.0.0 --port ${PORT}