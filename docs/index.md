# Mongoke documentation

Mongoke generates a graphql server based on a configuration file that describes the shape of the database via types and their corresponding collections.

Every type defined in the schema must be associated with a collection to be accessible via graphql, every type has a configuration to specify its collection and optionally authorization guards and disambiguations in case the type is an union or interface.

Types can also specify additional fields to connect to other entities via relations, this can be done in the `relations` part of the configuration.


- [Configuration](./configuration.md)
- [Graphql Queries](./queries.md)
- [Docker Usage and Env Vars](./docker.md)
