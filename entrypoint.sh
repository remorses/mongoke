#!/bin/sh
# pwd
python -m mongoke $1 && uvicorn generated.main:app --port 80