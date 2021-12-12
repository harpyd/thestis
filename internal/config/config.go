package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

const (
	defaultHTTPPort      = 8080
	defaultHTTPRWTimeout = 10 * time.Second
)

const LocalEnv = "local"

type (
	Config struct {
		Environment string
		HTTP        HTTP
	}

	HTTP struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}
)

func FromPath(configsPath string) (*Config, error) {
	setDefaults()

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = LocalEnv
	}

	if err := parseConfig(configsPath, appEnv); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.Environment = appEnv

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.readTimeout", defaultHTTPRWTimeout)
	viper.SetDefault("http.writeTimeout", defaultHTTPRWTimeout)
}

func parseConfig(configsPath, env string) error {
	viper.AddConfigPath(configsPath)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == LocalEnv {
		return replaceConfigEnvs()
	}

	viper.SetConfigName(env)

	if err := viper.MergeInConfig(); err != nil {
		return err
	}

	return replaceConfigEnvs()
}

func replaceConfigEnvs() error {
	var cmnErr error

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			envVal, err := envValue(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"))
			cmnErr = multierr.Append(cmnErr, err)

			value = envVal
		}

		viper.Set(k, value)
	}

	return cmnErr
}

func envValue(key string) (string, error) {
	envKey, defaultVal, hasDef := parseEnv(key)

	value, ok := os.LookupEnv(envKey)
	if !ok {
		if hasDef {
			return defaultVal, nil
		}

		return "", noEnvError{envKey: key}
	}

	return value, nil
}

func parseEnv(key string) (envKey, defaultValue string, hasDef bool) {
	s := strings.SplitN(key, ":", 2)
	envKey = s[0]

	if len(s) == 2 {
		defaultValue = s[1]
		hasDef = true
	}

	return envKey, defaultValue, hasDef
}

func unmarshal(cfg *Config) error {
	return viper.UnmarshalKey("http", &cfg.HTTP)
}

type noEnvError struct {
	envKey string
}

func IsNoEnvError(err error) bool {
	var enverr noEnvError

	return errors.As(err, &enverr)
}

func (e noEnvError) Error() string {
	return fmt.Sprintf("no %s env", e.envKey)
}
