## Environment variables
The accepted evn vars are:
- **DB_URL** the mongodb database url, **the url must contains the database name**, for example `mongodb://localhost/db`
- **DISABLE_GRAPHIQL** disable graphiql endpoint
- **GRAPHIQL_DEFAULT_QUERY** the query showed in the graphiql panel, default is ''
- **GRAPHIQL_DEFAULT_JWT** the default jwt used in graphiql headers panel

## Compose
The docker image is hosted on docker hub and has name `mongoke/mongoke`, every time the container is started the python code is generated based on the configuration given at path `/conf.yml`.
When used in docker-compose the code won't be generated unless the configuration changed, to force the regeneration use the command `docker-compose up -V`
When using an external skema configuration you must mount the file inside the container as a volume.
An example on a docker-compose file:
```yml
version: '3'

services:
    mongoke:
        image: mongoke/mongoke
        volumes: 
            - ./confs/testing.yaml:/conf.yml
            - ./confs/skema:/skema
        environment: 
            - DB_URL=mongodb://mongo/db
    mongo:
        image: mongo
        logging: 
            driver: none
```

## Swarm
When using with docker swarm probably you want to use configs instead of volumes to be able to deploy the containers in multiple nodes:
```yml
version: '3'
services:
    mongoke:
        image: mongoke/mongoke
        configs:
            - conf.yml
            - domain.skema
        environment: 
            - DB_URL=mongodb://mongo/db
    mongo:
        image: mongo
        logging: 
            driver: none

configs:
    domain.skema:
        name: domain-${domain_sum}.skema
        file: domain.skema
    conf.yml:
        name: mongoke_conf-${mongoke_sum}.yml
        file: mongoke_conf.yaml
```
Then to deploy the swarm add the env vars for the config versioning (configs in swarm are immutable and their name must change when the config content changes)
```sh
mongoke_sum=`cat mongoke_conf.yaml | md5` \
domain_sum=`cat domain.skema | md5` \
docker --host ssh://your-server.com stack deploy -c stack.yml stackname
```

## Apollo Federation
apollo federation permits you to gllue together different graphql servers together,
To add together more graphql servers together you can use the docker image [`xmorse/apollo-federation-gateway`](https://github.com/remorses/apollo-federation-gateway), this container acts as a gateway between your graphql servers and glue together the different schemas.
This approach is perfect to add mutations to your existing mongoke service.
Keep in mind that when you change one service you have also restart the gateway to make the changes to the federated schema, this can be handled adding a `POLL_INTERVAL` env variable that makes the gateway scrape the schemas every tot seconds to check if some schema changed.
An example using docker-compose:
```yml
# docker-compose.yml
version: '3'
services:
    gateway:
        image: xmorse/apollo-federation-gateway
        ports:
            - 8090:80
        environment:
            URL_0: "http://mongoke/"
            URL_1: "http://server"
    mongoke:
        image: mongoke/mongoke
        volumes: 
            - ./confs/testing.yaml:/conf.yml
            - ./confs/skema:/skema
        environment: 
            - DB_URL=mongodb://mongo/db
    server:
        build: server # your custom graphql server to handle mutations and stuff
    mongo:
        image: mongo
        logging: 
            driver: none
```

