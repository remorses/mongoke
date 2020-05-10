package mongoke

import (
	"github.com/PaesslerAG/gval"
	yaml "gopkg.in/yaml.v2"
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
}

type TypeConfig struct {
	Exposed     *bool       `yaml:"exposed"`
	Collection  string      `yaml:"collection"`
	Permissions []AuthGuard `yaml:"permissions"`
}

type AuthGuard struct {
	Expression        string   `yaml:"if"`
	AllowedOperations []string `yaml:"actions"`
	HideFields        []string `yaml:"hide_fields"`
	eval              gval.Evaluable
}

type RelationConfig struct {
	From         string                 `yaml:"from"`
	To           string                 `yaml:"to"`
	RelationType string                 `yaml:"relation_type"`
	where        map[string]interface{} `yaml:"where"`
}

// MakeConfigFromYaml parses the config from yaml
func MakeConfigFromYaml(data string) (Config, error) {
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
