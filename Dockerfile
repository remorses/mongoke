FROM python:3.7.2

RUN apt-get update && apt-get install -y dumb-init cmake bison flex git jq

WORKDIR /src

COPY requirements.txt /src/

RUN pip install -r requirements.txt

COPY mongoke /src/mongoke

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
# RUN python -m src 

ENTRYPOINT ["dumb-init", "--", "/entrypoint.sh"]
CMD ["/conf.yml"]

