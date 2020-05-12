package mongoke

import (
	"context"

	"github.com/PaesslerAG/gval"
	"github.com/caarlos0/env"
	yaml "github.com/ghodss/yaml"
)

type Map map[string]interface{}

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

const (
	ASC  = 1
	DESC = -1
)

var (
	True  = true
	False = false
)

type Config struct {
	DatabaseUri       string                 `json:"database_uri" env:"DB_URL"`
	DisableGraphiql   bool                   `json:"disable_graphiql" env:"DISABLE_GRAPHIQL"`
	Schema            string                 `json:"schema"`
	SchemaPath        string                 `json:"schema_path"`
	Types             map[string]*TypeConfig `json:"types"`
	Relations         []RelationConfig       `json:"relations"`
	JwtConfig         JwtConfig              `json:"jwt"`
	DatabaseFunctions DatabaseInterface
	Cache             Map
}

func (config Config) GetTypeConfig(typeName string) *TypeConfig {
	for name, conf := range config.Types {
		if name == typeName {
			return conf
		}
	}
	return nil
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

func (guard AuthGuard) Evaluate(params Map) (interface{}, error) {
	if guard.eval == nil {
		eval, err := gval.Full().NewEvaluable(guard.Expression)
		if err != nil {
			return nil, err
		}
		guard.eval = eval
	}
	res, err := guard.eval(context.Background(), params)
	if err != nil {
		return nil, err
	}
	return res, nil
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

	config := Config{}

	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		return Config{}, err
	}

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
