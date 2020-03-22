---
route: /docs
name: Quick Start
---

## Introduction

Mongoke generates a graphql server based on a configuration file that describes the shape of the database via graphql types and their corresponding collections.

Every type defined in the schema must be associated with a collection to be accessible via graphql, every type has a configuration to specify its collection and optionally authorization guards.

## Using Docker Compose

The fastest way to try Mongoke is via docker-compose.

### 1. Write the configuration to describe the database schema and relations

The ObjectId scalar is already defined by default, it is converted to string when sent as json

``` yaml
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

To start the container mount copy paste the following content in a `docker-compose.yml` file, then execute `docker-compose up` .

``` yaml
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

``` graphql
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

