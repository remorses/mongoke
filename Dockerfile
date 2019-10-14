FROM python:3.7.4-alpine

RUN apk update && apk add --no-cache build-base libffi-dev dumb-init cmake bison flex

WORKDIR /src

COPY requirements.txt /src/

RUN pip install -r requirements.txt

COPY mongoke /src/mongoke

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
# RUN python -m src 
RUN pip show tartiflette

ENTRYPOINT ["dumb-init", "--", "/entrypoint.sh"]
CMD ["/conf.yml"]

