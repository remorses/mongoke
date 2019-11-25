FROM tiangolo/uvicorn-gunicorn:python3.6-alpine3.8

RUN apk update && apk add --no-cache build-base libffi-dev dumb-init cmake bison flex

WORKDIR /src

COPY *.txt /src/

RUN pip install -r requirements.txt
RUN pip install -r requirements-generated.txt

COPY mongoke /src/mongoke

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
# RUN python -m src 
# RUN pip show tartiflette

ENV PYTHONUNBUFFERED=1

ENTRYPOINT ["dumb-init", "--", "/entrypoint.sh"]
CMD ["/conf.yml"]

