package config

import (
	"os"
	"path/filepath"
	"runtime"

	coreConfig "github.com/pixality-inc/golang-core/config"
	"github.com/pixality-inc/golang-core/http"
	"github.com/pixality-inc/golang-core/logger"
	"github.com/pixality-inc/golang-core/postgres"
)

type Config struct {
	Logger      logger.YamlConfig           `env-prefix:"PIXALITY_LOG_"         yaml:"logger"`
	Http        http.ConfigYaml             `env-prefix:"PIXALITY_HTTP_"        yaml:"http"`
	AdminHttp   http.ConfigYaml             `env-prefix:"PIXALITY_ADMIN_HTTP_"  yaml:"admin_http"`
	Healthcheck http.ConfigYaml             `env-prefix:"PIXALITY_HEALTHCHECK_" yaml:"healthcheck"`
	Metrics     http.ConfigYaml             `env-prefix:"PIXALITY_METRICS_"     yaml:"metrics"`
	Database    postgres.DatabaseConfigYaml `env-prefix:"PIXALITY_DB_"          yaml:"database"`
}

func RootDir() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Join(filepath.Dir(b), "../..")
	)

	return basepath
}

func configFile() string {
	configFilename := os.Getenv("PIXALITY_CONFIG_FILE")
	if configFilename == "" {
		configFilename = filepath.Join(RootDir(), "config.yaml")
	}

	return configFilename
}

func LoadConfig() *Config {
	return coreConfig.LoadConfig[Config](configFile())
}
