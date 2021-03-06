package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

type Env = string

const LocalEnv Env = "local"

type StepsPolicy = string

const SavePerStepPolicy StepsPolicy = "savePerStep"

type SignalBus = string

const Nats SignalBus = "nats"

type AuthType = string

type LoggerLib = string

const Zap LoggerLib = "zap"

type LoggerLevel = string

const (
	DebugLevel LoggerLevel = "DEBUG"
	InfoLevel  LoggerLevel = "INFO"
	WarnLevel  LoggerLevel = "WARN"
	ErrorLevel LoggerLevel = "ERROR"
	FatalLevel LoggerLevel = "FATAL"
)

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
		Pipeline    Pipeline
		SavePerStep SavePerStep
		Nats        NatsServer
		Logger      Logger
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

	Pipeline struct {
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

	Logger struct {
		Lib   LoggerLib
		Level LoggerLevel
	}
)

const (
	defaultHTTPPort            = 8080
	defaultHTTPRWTimeout       = 10 * time.Second
	defaultHTTPShutdownTimeout = 10 * time.Second
)

const defaultMongoDisconnectTimeout = 10 * time.Second

const (
	defaultPipelineFlowTimeout = 24 * time.Hour
	defaultPipelineWorkers     = 100
)

const (
	defaultLoggerLib   = Zap
	defaultLoggerLevel = "INFO"
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
	viper.SetDefault("pipeline.flowTimeout", defaultPipelineFlowTimeout)
	viper.SetDefault("pipeline.workers", defaultPipelineWorkers)
	viper.SetDefault("logger.lib", defaultLoggerLib)
	viper.SetDefault("logger.level", defaultLoggerLevel)
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

	if err := viper.UnmarshalKey("pipeline", &cfg.Pipeline); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("savePerStep", &cfg.SavePerStep); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("nats", &cfg.Nats); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("firebase", &cfg.Firebase); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("logger", &cfg.Logger); err != nil {
		return err
	}

	cfg.Logger.Level = strings.ToUpper(cfg.Logger.Level)

	return nil
}

type noEnvError struct {
	envKey string
}

func (e noEnvError) Error() string {
	return fmt.Sprintf("no %s env", e.envKey)
}
