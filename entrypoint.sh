#!/bin/sh
# pwd
python -m mongoke $1 && uvicorn generated.main:app --host 0.0.0.0 --port 80