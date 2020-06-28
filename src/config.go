package goke

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PaesslerAG/gval"
	"github.com/caarlos0/env"
	yaml "github.com/ghodss/yaml"
)

type Map = map[string]interface{}

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

var AdminSecretHeader = "Goke-Admin-Header"

var DEFAULT_PERMISSIONS = []string{
	Operations.READ,
	Operations.CREATE,
	Operations.UPDATE,
	Operations.DELETE,
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
	// DatabaseUri       string
	Admins             []AdminConfig          `json:"admins"`
	Mongodb            MongodbConfig          `json:"mongodb"`
	FakeDatabase       FakeDatabaseConfig     `json:"fake_database"`
	Firestore          FirestoreConfig        `json:"firestore"`
	DisableGraphiql    bool                   `json:"disable_graphiql" env:"DISABLE_GRAPHIQL"`
	Schema             string                 `json:"schema"`
	SchemaPath         string                 `json:"schema_path"`
	SchemaUrl          string                 `json:"schema_url"`
	Types              map[string]*TypeConfig `json:"types"`
	Relations          []RelationConfig       `json:"relations"`
	JwtConfig          JwtConfig              `json:"jwt"`
	DefaultPermissions []string               `json:"default_permissions"`
	DatabaseFunctions  DatabaseInterface
	Cache              Map
}

func (config *Config) Init() error {
	// types cache (gaphql-go complains of duplicate types)
	if config.Cache == nil {
		config.Cache = make(Map)
	}

	// add schema type defs
	if config.Schema == "" && config.SchemaPath != "" {
		data, e := ioutil.ReadFile(config.SchemaPath)
		if e != nil {
			return e
		}
		config.Schema = string(data)
	}
	if config.Schema == "" && config.SchemaUrl != "" {
		data, e := DownloadFile(config.SchemaUrl)
		if e != nil {
			return e
		}
		config.Schema = string(data)
	}
	if config.Schema == "" {
		return errors.New("missing required schema")
	}

	// default permissions
	if config.DefaultPermissions == nil {
		config.DefaultPermissions = DEFAULT_PERMISSIONS
	}

	// initialize the expressions
	for _, t := range config.Types {
		for _, guard := range t.Permissions {
			err := guard.Init()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type AdminConfig struct {
	Secret string `json:"secret"`
}

type MongodbConfig struct {
	Uri string `json:"uri" env:"MONGODB_URL"` // TODO does env parse create these structs?
}

type FakeDatabaseConfig struct {
	DocumentsPerCollection *int `json:"documents_per_collection"`
}

type FirestoreConfig struct {
	ProjectID string `json:"project_id" env:"FIRESTORE_PROJECT_ID"`
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
	AllowedOperations []string `json:"allow_operations"`
	// HideFields        []string `json:"hide_fields"`
	eval gval.Evaluable
}

func (guard *AuthGuard) Init() error {
	// println("initializing guard with expression " + guard.Expression)
	eval, err := gval.Full().NewEvaluable(guard.Expression)
	if err != nil {
		return err
	}
	guard.eval = eval
	return nil
}

func (guard *AuthGuard) Evaluate(params Map) (interface{}, error) {
	fmt.Printf("evaluating '%s'", guard.Expression)
	if guard.eval == nil {
		guard.Init()
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
	// TODO currently the relations do not support and, or operators
}

func InterpolateMatch(match map[string]Filter, scope Map) (map[string]Filter, error) {
	result := make(map[string]Filter)
	for k, filter := range match {
		interpolated, err := filter.Interpolate(scope)
		if err != nil {
			return nil, err
		}
		result[k] = interpolated
	}
	return result, nil
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(url string) (string, error) {
	// TODO i should test downloadFile
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	s := buf.String() // Does a complete copy of the bytes in the buffer.

	return s, err
}

func ReverseMaps(ss []Map) []Map {
	if len(ss) == 0 {
		return ss
	}
	copy := make([]Map, len(ss))
	j := 0
	for i := len(ss) - 1; i >= 0; i-- {
		copy[j] = ss[i]
		j++
	}
	return copy
}
