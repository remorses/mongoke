# cli

problem this solves

-   no docker
-   people can't find the url of the endpoint
-   easy mocks
-   easy watching
-   easy typescript graphql client

mongoke -p port --mock -w

-   starts the server on port 8030 by default
-   generates the graphql typed client generateClient.path directory
-   starts the mock server if --mock
-   if -w watch the configuration and mocks file for changes

## start the server

generate the server code in a .mongoke directory
run the server via python3 and uvicorn

## mocking

start the mongoke server at a random port, then start the mock server pointing it to the mongoke server, this fetches the schema and exposes a new graphql endpoint

Problem arises if i expose more than a graphiql in that
