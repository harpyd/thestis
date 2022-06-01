package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

const (
	defaultHTTPPort               = 8080
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPShutdownTimeout    = 10 * time.Second
	defaultMongoDisconnectTimeout = 10 * time.Second
	defaultPerformanceFlowTimeout = 24 * time.Hour
	defaultPerformanceWorkers     = 100
)

type Env = string

const LocalEnv Env = "local"

type StepsPolicy = string

const SavePerStepPolicy StepsPolicy = "savePerStep"

type SignalBus = string

const Nats SignalBus = "nats"

type AuthType = string

const (
	FakeAuth     AuthType = "fake"
	FirebaseAuth AuthType = "firebase"
)

type (
	Config struct {
		Environment Env
		HTTP        HTTP
		Mongo       Mongo
		Auth        Auth
		Firebase    Firebase
		Performance Performance
		SavePerStep SavePerStep
		Nats        NatsServer
	}

	HTTP struct {
		Port            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		ShutdownTimeout time.Duration
		AllowedOrigins  []string
	}

	Mongo struct {
		URI               string
		DatabaseName      string
		Username          string
		Password          string
		DisconnectTimeout time.Duration
	}

	Auth struct {
		With string
	}

	Firebase struct {
		ServiceAccountFile string
	}

	Performance struct {
		FlowTimeout time.Duration
		Policy      StepsPolicy
		SignalBus   SignalBus
		Workers     int
	}

	SavePerStep struct {
		SaveTimeout time.Duration
	}

	NatsServer struct {
		URL string
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
	viper.SetDefault("http.shutdownTimeout", defaultHTTPShutdownTimeout)
	viper.SetDefault("mongo.disconnectTimeout", defaultMongoDisconnectTimeout)
	viper.SetDefault("performance.flowTimeout", defaultPerformanceFlowTimeout)
	viper.SetDefault("performance.workers", defaultPerformanceWorkers)
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
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("performance", &cfg.Performance); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("savePerStep", &cfg.SavePerStep); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("nats", &cfg.Nats); err != nil {
		return err
	}

	return viper.UnmarshalKey("firebase", &cfg.Firebase)
}

type noEnvError struct {
	envKey string
}

func (e noEnvError) Error() string {
	return fmt.Sprintf("no %s env", e.envKey)
}
