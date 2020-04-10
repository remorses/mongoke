# examples

To try an example go inside its directory and execute `docker-compose up`, all the examples expose the graphiql at [`http://localhost:8090`](http://localhost:8090) and set some default queries for you to try, these queries are inside the `quieries.grpaphql` files.
Before try other examples remember to remove the previous containers with `docker.compose down` or `docker kill $(docker ps -q)`, or the port `8090` will be occupied by the previous containers.

-   [Basics](./basic)
-   [Realtions example](./relations)
-   [Union types](./unions)
-   ~~[Interface types](./interfaces)~~
-   [Authorization Guards](./authorizartion)
-   ~~[Using skema to speed up other services developement]()~~
-   [Usage with apollo federation](./federation)
-   [Real world usage](./real_world)
