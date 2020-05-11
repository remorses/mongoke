package mongoke

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func validate(config Config) {

	schemaLoader := gojsonschema.NewStringLoader(jsonSchemaString)
	result, err := gojsonschema.Validate(schemaLoader, nil)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}

var jsonSchemaString = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$ref": "#/definitions/Configuration",
    "definitions": {
        "Configuration": {
        "type": "object",
            "required": [
                "types"
            ],
            "properties": {
                "schema": {
                    "type": "string",
                },
                "schema_url": {
                    "$ref": "#/definitions/Url",
                },
                "schema_path": {
                    "type": "string",
                },
                "types": {
                "description": "",
                    "type": "object",
                    "required": [],
                    "properties": {},
                    "additionalProperties": {
                        "description": "",
                        "type": "object",
                        "required": [
                            "collection"
                        ],
                        "properties": {
                            "collection": {
                                "type": "string",
                            },
                            "exposed": {
                                "type": "boolean",
                            },
                            "pipeline": {
                                "type": "array",
                            "items": {
                                },
                                "minItems": 0
                            },
                            "disambiguations": {
                            "description": "",
                                "type": "object",
                                "required": [],
                                "properties": {},
                                "additionalProperties": {
                                    "type": "string",
                                }
                            },
                            "guards": {
                                "type": "array",
                            "items": {
                                "description": "",
                                    "type": "object",
                                    "required": [
                                        "expression"
                                    ],
                                    "properties": {
                                        "expression": {
                                            "type": "string",
                                        },
                                        "excluded": {
                                            "type": "array",
                                        "items": {
                                                "type": "string",
                                            },
                                            "minItems": 0
                                        },
                                        "when": {
                                            "enum": [
                                                "after",
                                                "before"
                                            ],
                                            "type": "string",
                                        }
                                    },
                                    "additionalProperties": false
                                },
                                "minItems": 0
                            }
                        },
                        "additionalProperties": false
                    }
                },
                "relations": {
                    "type": "array",
                "items": {
                    "description": "",
                        "type": "object",
                        "required": [
                            "from",
                            "to",
                            "relation_type",
                            "field",
                            "where"
                        ],
                        "properties": {
                            "from": {
                                "type": "string",
                            },
                            "to": {
                                "type": "string",
                            },
                            "relation_type": {
                                "enum": [
                                    "to_many",
                                    "to_one"
                                ],
                                "type": "string",
                            },
                            "field": {
                                "type": "string",
                            },
                            "where": {
                            }
                        },
                        "additionalProperties": false
                    },
                    "minItems": 0
                },
                "jwt": {
                "description": "",
                    "type": "object",
                    "required": [],
                    "properties": {
                        "secret": {
                            "type": "string",
                        },
                        "header_name": {
                            "type": "string",
                        },
                        "header_scheme": {
                            "type": "string",
                        },
                        "required": {
                            "type": "boolean",
                        },
                        "algorithms": {
                            "type": "array",
                        "items": {
                                "enum": [
                                    "H256",
                                    "HS512",
                                    "HS384",
                                    "RS256",
                                    "RS384",
                                    "RS512",
                                    "ES256",
                                    "ES384",
                                    "ES521",
                                    "ES512",
                                    "PS256",
                                    "PS384",
                                    "PS512"
                                ],
                                "type": "string",
                            },
                            "minItems": 0
                        }
                    },
                    "additionalProperties": false
                }
            },
            "additionalProperties": false
        },
        "Url": {
            "type": "string",
        }
    }
}
`
