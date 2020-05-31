package mongoke

import (
	"bytes"
	"context"
	"net/http"

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
	DatabaseUri       string
	Mongodb           MongodbConfig          `json:"mongodb"`
	FakeDatabase      FakeDatabaseConfig     `json:"fake_database"`
	Firestore         FirestoreConfig        `json:"firestore"`
	DisableGraphiql   bool                   `json:"disable_graphiql" env:"DISABLE_GRAPHIQL"`
	Schema            string                 `json:"schema"`
	SchemaPath        string                 `json:"schema_path"`
	SchemaUrl         string                 `json:"schema_url"`
	Types             map[string]*TypeConfig `json:"types"`
	Relations         []RelationConfig       `json:"relations"`
	JwtConfig         JwtConfig              `json:"jwt"`
	DatabaseFunctions DatabaseInterface
	Cache             Map
}

type MongodbConfig struct {
	Uri string `json:"uri" env:"MONGODB_URL"`
}

type FakeDatabaseConfig struct {
	DocumentsPerCollection int `json: "documents_per_collection"`
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
