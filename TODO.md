## Todo:

-   ~~add a default graphql text to show how mongoke queries work~~
-   add an introduction on how to install docker compose
-   remove plurals and use the Nodes postfix instead

-   ~~exceptions are returned as `Server got an error`, use an error with `coerce_value` to return e better message~~
-   ~~empty variables search gets nothing~~
-   option to allow only some origins in cors, default is \*
-   when a type is in a relation but not in the types section, KeyError is thrown at launch
-   ~~Enum and Scalar values are not searchable in where~~
-   ~~dont lowercase the first letter of query fields~~
-   make smaller docker image
-   ~~remove hard limit of resolved nodes, pass the limit in the configuration~~
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
-   ~~better performance of connection_resolver removing the $skip and $count~~
-   add a dataloader for single connections
