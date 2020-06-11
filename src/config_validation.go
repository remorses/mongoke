package goke

import (
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

func validateYamlConfig(data string) error {
	var document map[string]interface{}
	if err := yaml.Unmarshal([]byte(data), &document); err != nil {
		return err
	}

	schemaLoader := gojsonschema.NewStringLoader(jsonSchemaString)
	result, err := gojsonschema.Validate(schemaLoader, gojsonschema.NewGoLoader(document))
	if err != nil {
		return err
	}

	if !result.Valid() {
		msg := "The yaml config is not valid. see errors:\n"
		for _, desc := range result.Errors() {
			msg += fmt.Sprintf("- %s\n", desc)
		}
		return errors.New(msg)
	}
	return nil
}

// TODO upload jsonschema somewhere
var jsonSchemaString = `
{
    "$schema": "http://json-schema.org/draft-07/schema#",
    "$ref": "#/definitions/Configuration",
    "definitions": {
        "Configuration": {
            "type": "object",
			"required": ["types"],
			"anyOf": [
				{"required": ["schema"]},
				{"required": ["schema_path"]},
				{"required": ["schema_url"]}
			],
            "properties": {
                "schema": {
                    "type": "string"
                },
				"mongodb": {
					"type": "object",
					"properties": {
						"uri": {
							"type": "string"
						}
					}
				},
				"fake_database": {
					"type": "object",
					"properties": {
						"documents_per_collection": {
							"type": "number"
						}
					}
				},
				"firestore": {
					"type": "object",
					"properties": {
						"uri": {
							"type": "string"
						}
					}
				},
				"default_permissions": {
					"type": "array",
                    "items": {
                        "$ref": "#/definitions/PermissionEnum"
                    }
				},
                "schema_url": {
                    "$ref": "#/definitions/Url"
                },
                "schema_path": {
                    "type": "string"
                },
                "types": {
					"type": "object",
					"minProperties": 1,
                    "properties": {},
                    "additionalProperties": {
                        "type": "object",
                        "required": ["collection"],
                        "properties": {
                            "collection": {
                                "type": "string"
                            },
                            "exposed": {
                                "type": "boolean"
                            },
                            "type_check": {
                                "type": "string"
                            },
                            "permissions": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/AuthGuard"
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
                        "type": "object",
                        "required": [
                            "from",
                            "to",
                            "type",
                            "field",
                            "where"
                        ],
                        "properties": {
                            "from": {
                                "type": "string"
                            },
                            "to": {
                                "type": "string"
                            },
                            "type": {
                                "enum": ["to_many", "to_one"],
                                "type": "string"
                            },
                            "field": {
                                "type": "string"
                            },
                            "where": {
                                "$ref": "#/definitions/WhereFilter"
                            }
                        },
                        "additionalProperties": false
                    },
                    "minItems": 0
                },
                "jwt": {
                    "type": "object",
                    "properties": {
                        "key": {
                            "type": "string"
                        },
                        "jwk_url": {
                            "type": "string"
                        },
                        "header_name": {
                            "type": "string"
                        },
                        "audience": {
                            "type": "boolean"
                        },
                        "issuer": {
                            "type": "boolean"
                        },
                        "type": {
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
                            "type": "string"
                        }
                    },
                    "additionalProperties": false
                }
            },
            "additionalProperties": false
        },
        "Url": {
            "type": "string"
        },
        "WhereFilter": {
            "type": "object",
            "additionalProperties": {
                "eq": {},
                "neq": {},
                "in": {
                    "type": "array",
                    "items": {}
                },
                "nin": {
                    "type": "array",
                    "items": {}
                }
            }
        },
        "AuthGuard": {
            "type": "object",
            "required": ["if"],
            "properties": {
                "if": {
                    "type": "string"
                },
                "allow_operations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/PermissionEnum"
                    },
                    "minItems": 0
                }
            },
            "additionalProperties": false
        },
        "PermissionEnum": {
            "enum": ["read", "update", "delete", "create"],
            "type": "string"
        }
    }
}
`
