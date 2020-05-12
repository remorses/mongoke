package mongoke

import (
	"github.com/PaesslerAG/gval"
	yaml "github.com/ghodss/yaml"
)

var Operations = struct {
	READ   string
	UPDATE string
	DELETE string
	CREATE string
}{
	READ:   "read",
	UPDATE: "update",
	DELETE: "delete",
	CREATE: "create",
}

type Config struct {
	EnableGraphiql    *bool                  `json:"enable_graphiql"`
	DatabaseUri       string                 `json:"database_uri"`
	Schema            string                 `json:"schema"`
	SchemaPath        string                 `json:"schema_path"`
	Types             map[string]*TypeConfig `json:"types"`
	Relations         []RelationConfig       `json:"relations"`
	JwtConfig         JwtConfig              `json:"jwt"`
	databaseFunctions DatabaseInterface
	cache             Map
}

type JwtConfig struct {
	HeaderName string `json:"header_name"`
	// for example https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com
	JwkUrl string `json:"jwk_url"`
	// 32 string for HS256 or PEM encoded string or as a X509 certificate for RSA
	Key string `json:"key"`
	// HS256, HS384, HS512, RS256, RS384, RS512
	Type     string `json:"type"`
	Audience string `json:"audience"`
	Issuer   string `json:"issuer"`
}

type TypeConfig struct {
	Exposed     *bool       `json:"exposed"`
	Collection  string      `json:"collection"`
	IsTypeOf    string      `json:"type_check"`
	Permissions []AuthGuard `json:"permissions"`
}

type AuthGuard struct {
	Expression        string   `json:"if"`
	AllowedOperations []string `json:"allowed_operations"`
	HideFields        []string `json:"hide_fields"`
	eval              gval.Evaluable
}

type RelationConfig struct {
	From         string            `json:"from"`
	Field        string            `json:"field"`
	To           string            `json:"to"`
	RelationType string            `json:"type"`
	Where        map[string]Filter `json:"where"`
}

// MakeConfigFromYaml parses the config from yaml
func MakeConfigFromYaml(data string) (Config, error) {

	if err := validateYamlConfig(data); err != nil {
		return Config{}, err
	}

	t := Config{}

	// TODO add databaseUri overwite from environment

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return Config{}, err
	}

	return t, nil
}

func (config Config) getTypeConfig(typeName string) *TypeConfig {
	for name, conf := range config.Types {
		if name == typeName {
			return conf
		}
	}
	return nil
}
