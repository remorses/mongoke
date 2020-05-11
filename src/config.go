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
	DatabaseUri       string                 `yaml:"database_uri"`
	Schema            string                 `yaml:"schema"`
	SchemaPath        string                 `yaml:"schema_path"`
	Types             map[string]*TypeConfig `yaml:"types"`
	Relations         []RelationConfig       `yaml:"relations"`
	JwtConfig         JwtConfig              `yaml:"jwt"`
	databaseFunctions DatabaseInterface
}

type JwtConfig struct {
	HeaderName string `yaml:"header_name"`
	// for example https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com
	JwkUrl string `yaml:"jwk_url"`
	// 32 string for HS256 or PEM encoded string or as a X509 certificate for RSA
	Key string `yaml:"key"`
	// HS256, HS384, HS512, RS256, RS384, RS512
	Type     string `yaml:"type"`
	Audience string `yaml:"audience"`
	Issuer   string `yaml:"issuer"`
}

type TypeConfig struct {
	Exposed     *bool       `yaml:"exposed"`
	Collection  string      `yaml:"collection"`
	IsTypeOf    string      `yaml:"type_check"`
	Permissions []AuthGuard `yaml:"permissions"`
}

type AuthGuard struct {
	Expression        string   `yaml:"if"`
	AllowedOperations []string `yaml:"allowed_operations"`
	HideFields        []string `yaml:"hide_fields"`
	eval              gval.Evaluable
}

type RelationConfig struct {
	From         string            `yaml:"from"`
	Field        string            `yaml:"field"`
	To           string            `yaml:"to"`
	RelationType string            `yaml:"type"`
	Where        map[string]Filter `yaml:"where"`
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
