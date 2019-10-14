FROM python:3.7.4-alpine as base

RUN apk update && apk add build-base libffi-dev dumb-init cmake bison flex

WORKDIR /src

COPY requirements.txt /src/

RUN pip install --user -r requirements.txt

FROM python:3.7.4-alpine

RUN apk add dumb-init

COPY --from=base /root/.local /root/.local

WORKDIR /src

COPY mongoke /src/mongoke

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["dumb-init", "--", "/entrypoint.sh"]
CMD ["/conf.yml"]

