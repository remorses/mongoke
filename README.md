Jsonschema to graphql mongodb server


cose da aggiungere:
- devo determinare il _typename di alcuni oggetti, questo si può ottenere immettendo un _typename apposta nel database, oppure aggiungendo un file di configurazione in più
- devo creare dei types in più per ogni collection, per gestire le connections, gli where, orderBy, 
- devo aggiungere un file di configurazione per gestire le relations
- le autorizzazioni sono gestite via configurazione apposita, per ogni type viene associata una rule, la rule è una espressione in python che viene eseguita prima di concedere la risorsa, ha accesso agli headers, field della risorsa. quando un utente cerca di fare una query ad una risorsa, viene prima eseguita una verifica, la espressione viene eseguita sui campi del where e poi sui campi della risorsa vera e propria, se l’espressione ritorna false accesso negato, se il where non contiene parti della espressione allora l’accesso è negato

Fasi:
- genera i types principali attraverso skema
- Aggiungi gli scalars come objectid, ...
- Genera il type Query insieme a tutto il boilerplate
- Aggiungi le relation fields usando extend types, aggiungendo altro boilerplate
- Genera i resolver per la query, attraverso un file di template
- Genera i resolver per le relations, applicando il where preso dalla configurazione
- Per gli Union types uso disambiguazioni per aggiungere typename alla fine


Create the base graphql types in graphql,
Then generate the basic queries to get the collection directly via graphql, 
Like
users(where, orderBy, first, last, after, before)
user(where)

Le relations sono descritte in un file di configurazione a parte,
Per ogni relation viene aggiunto un field al graphql schema e il corrispondente resolver, basato su un field specifico del parent

Per esempio
bots: Campaign.bots_ids -> Bot

In questo modo verrà aggiunto il fields bots al type campaign

Il where delle queries viene passato così come è a mongodb,
lo stesso per OrderBy
After e before sono basati sul valore di orderBy, quindi after: x con orderBy: author significa applicare author: gte: x

Dovrò aggiungere un layer di batch per evitare di fare troppe richieste al database, questo si può ottenere utilizzando un dataloader


—


to generate the schema

every type gets:

```graphql

type Query {
    ${type}s(where: Where${type}, orderBy: OrderBy,): ${type}Connection
    ${type}(where: Where${type}, ): ${type}
}

type ${type}Connection {
        nodes: [${type}]
        pageInfo: PageInfo
}

input Where${type} { 
    ${field}: GeneralWhere
}

input GeneralWhere {
    $in: [string]
    $eq: String
    $gte: NumberOrString
}

input OrderBy {
    ${field}: Direction
}

enum Direction {
    ASC
    DESC
}

# relations

extend ${fromType} {
    ${relationName}: ${toType} # if one_to_one
    ${relationName}(where, orderBy): Connection${toType} # if one_to_many
    ${relationName}(where, orderBy): Connection${toType} # if many_to_many
}

scalar NumberOrString
```