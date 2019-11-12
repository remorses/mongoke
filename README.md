<p align="center">
  <img width="300" src="https://github.com/remorses/mongoke/blob/master/.github/logo.jpg?raw=true">
</p>
<h1 align="center">mongoke</h1>
<h3 align="center">Instantly serve your MongoDb database via graphql</h3>

[**Docs**](https://remorses.github.io/mongoke/docs) • [**Examples**](https://github.com/remorses/mongoke-examples)

## Features

-   **Powerful Queries**: Pagination, filtering, relation, relay-style connections built-in and generated in a bunch of seconds
-   **Works with existing databases**: Point it to an existing MongoDb database to instantly get a ready-to-use GraphQL API
-   **Authorization via Jwt**: Every collection can be protected based on jwt payload and document fields
-   **Horizontally Scalable**: The service is completely stateless and can be replicated on demand
-   **Apollo Federation**: The service can be easily glued with other graphql servers to handle writes and more complicated logic.
-   **Resilient Idempotent Configuration**: One YAML Configuration as the only source of truth, relations, authorization and types in one file

## Quickstart:

### Docker compose

The fastest way to try Mongoke is via docker-compose.

### 1. Write the configuration to describe the database schema and relations

The ObjectId scalar is already defined by default, it is converted to string when sent as json

```yml
# ./mongoke.yml
schema: |
    type User {
        _id: ObjectId
        username: String
        email: String
    }
    type BlogPost {
        _id: ObjectId
        author_id: ObjectId
        title: String
        content: String
    }

types:
    User:
        collection: users
    BlogPost:
        collection: posts

relations:
    - field: posts
      from: User
      to: BlogPost
      relation_type: to_many
      where:
          author_id: ${{ parent['_id'] }}
```

### 2. Run the `mongoke` image with the above configuration

To start the container mount copy paste the following content in a `docker-compose.yml` file, then execute `docker-compose up`.

```yml
# docker-compose.yml
version: '3'

services:
    mongoke:
        ports:
            - 4000:80
        image: mongoke/mongoke
        environment:
            DB_URL: mongodb://mongo/db
        volumes:
            - ./mongoke.yml:/conf.yml
    mongo:
        image: mongo
        logging:
            driver: none
```

### 3. Query the generated service via graphql or go to [http://localhost:4000/graphiql](http://localhost:4000/graphiql) to open graphiql

```graphql
{
    user(where: { username: { eq: "Mike" } }) {
        _id
        username
        email
        posts {
            nodes {
                title
            }
        }
    }

    blogPosts(first: 10, after: "Post 1", cursorField: title) {
        nodes {
            title
            content
        }
        pageInfo {
            endCursor
            hasNextPage
        }
    }
}
```

---

## Tutorials

Please help the project making new tutorials and submit a issue to list it here!

---

## Todo:

-   ~~use graphql to define the schema~~
-   ~~publish the docker image (after tartiflette devs fix extend type issue)~~
-   ~~resolve issue connection nodes must all have an \_id field because it is default cursor field~~
-   integration tests for all the resolver types
-   integration tests for the relations
-   ~~add pipelines feature to all resolvers (adding a custom find and find_one made with aggregate)~~
-   ~~add the $ to the where input fields inside resolvers (in must be $in, ...)~~
-   ~~remove strip_nones after asserting v1 works~~

Low priority

-   ~~`required` config field, add verify the jwt with the secret if provided~~
-   ~~add schema validation to the configuration~~
-   subscriptions
-   ~~add `edges` to make connection type be relay compliant~~
-   better performance of connection_resolver removing the $skip and $count
-   add a dataloader for single connections
