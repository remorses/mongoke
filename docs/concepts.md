## Architecture

Mongoke generates a graphql server based on a configuration file that tells the shape of the database via types and their corresponding collections.
The database schema is defined via the skema language, an sdl that compiles to graphql types among many other languages.
Read more about the types here.

Every type defined in the schema must be associated with a collection to be accessible via graphql, every type has a configuration to specify its collection and optionally authorization guards and disambiguations in case the type is an union or interface.
Read more about authorization here and about disabmiguation here.

Types can also specify additional fields to connect to other entities via relations, this can be done in the `relations` part of the configuration.
Read more about relations here.




