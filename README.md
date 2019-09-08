Todo:
- cursor must be obfuscated in connection, (also after and before are string so it is a must)
- add pipelines feature to all resolvers (adding a custom find and find_one made with aggregate)
- ~~add the $ to the where input fields inside resolvers (in must be $in, ...)~~
- ~~remove strip_nones after asserting v1 works~~

Low priority
- add verify the jwt with the secret if provided
- add schema validation to the configuration
- add subscriptions
- add edges to make connection type be relay compliant 
- better performance of connection_resolver removing the $skip and $count
- add a dataloader for single connections