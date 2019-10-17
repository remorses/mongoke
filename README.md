<p align="center">
  <img width="300" src="https://github.com/remorses/mongoke/blob/master/.github/logo.jpg?raw=true">
</p>
<h1 align="center">mongoke</h1>
<h3 align="center">Instantly serve your MongoDb database via graphql<h3>
supports 
- authorization of the fields based on jwt, 
- relations between collections,
- mongodb aggregated pipelines
- apollo federation

## Usage
Mongoke serve your mongodb database via a declarative, idempotent configuration that describes the shape of the types in the database and their relations.
To get started first describe the shape of your types inside the database via the [skema](https://github.com/remorses/skema) language, then write a configuration for every type to connect it to the associated collection and add authorization guards.
Then you can add relations between types, describing what field will lead to the related types and if the relation is of type `to_one` or `to_many`.


```yml
# docker-compose.yml
version: '3'

services:
    mongoke:
        image: mongoke/mongoke
        environment: 
            PYTHONUNBUFFERED: '1'
            DB_URL: mongodb://mongo:27017   
            CONFIG:
                skema: |
                    User:
                        _id: ObjectId
                        username: Str
                        email: Str
                    BlogPost:
                        _id: ObjectId
                        author_id: ObjectId
                        content: Str
                types:
                    User:
                        collection: users
                    BlogPost:
                        collection: posts
                relations:
                    -   field: posts
                        from: User
                        to: BlogPost
                        where:
                            author_id: ${{ parent['_id'] }}
    mongo:
        image: mongo
        logging: 
            driver: none

```

Here is an example:
```yaml
# example.yml

db_url: mongodb://mongo:27017/db

skema: |
    Article:
        content: Str
        autorId: ObjectId
        createdAt: DateTime
    User:
        _id: ObjectId
        name: Str
        surname: Str
        aricleIds: [ObjectId]
    ObjectId: Any
    DateTime: Any

types:
    User:
        collection: users
    Article:
        collection: articles

relations:
    -   from: User
        to: Article
        relation_type: to_many
        field: articles
        query:
            autorId: ${{ parent['_id'] }}
```

Then generate the server code and serve it with the mongoke docker image
```
version: '3'

services:
    server:
        image: mongoke/mongoke:latest
        command: /conf.yml
        volumes: 
            - ./example.yml:/conf.yml
    mongo:
        image: mongo
        logging: 
            driver: none
```

Then you can query the database from your graphql app as you like

```graphql

{
  author(where: {name: "Joseph"}) {
    name
    articles {
      nodes {
        content
      }
    }
  }
}
```

```graphql

{
  articles(first: 5, after: "22/09/1999", cursorField: createdAt) {
    nodes {
      content
    }
    pageInfo {
      endCursor
    }
  }
}
```

## Todo:
- ~~publish the docker image (after tartiflette devs fix extend type issue)~~
- resolve issue connection nodes must all have an _id field because it is default cursor field
- integration tests for all the resolver types
- integration tests for the relations
- cursor must be obfuscated in connection, (also after and before are string so it is a must)
- ~~add pipelines feature to all resolvers (adding a custom find and find_one made with aggregate)~~
- ~~add the $ to the where input fields inside resolvers (in must be $in, ...)~~
- ~~remove strip_nones after asserting v1 works~~

Low priority
- add verify the jwt with the secret if provided
- ~~add schema validation to the configuration~~
- add subscriptions
- add edges to make connection type be relay compliant 
- better performance of connection_resolver removing the $skip and $count
- add a dataloader for single connections
