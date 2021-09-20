package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
var (
	version    string = "dev"
	commitHash string = "-"
)

var cfg *Config

func GetConfigInstance() Config {
	if cfg != nil {
		return *cfg
	}

	return Config{}
}

// Database contains all parameters database connection.
type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslmode"`
	Driver   string `yaml:"driver"`
}

// GRPC contains parameter address gRPC.
type GRPC struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	MaxConnectionIdle int64  `yaml:"maxConnectionIdle"`
	Timeout           int64  `yaml:"timeout"`
	MaxConnectionAge  int64  `yaml:"maxConnectionAge"`
}

// REST contains parameter REST JSON connection.
type REST struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Project contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	SwaggerDir  string `yaml:"swaggerDir"`
	Version     string
	CommitHash  string
}

// Metrics contains all parameters metrics information.
type Metrics struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// Jaeger contains all parameters metrics information.
type Jaeger struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Service string `yaml:"service"`
}

// Kafka contains all parameters Kafka information.
type Kafka struct {
	Capacity uint64   `yaml:"capacity"`
	Topic    string   `yaml:"topic"`
	GroupID  string   `yaml:"groupId"`
	Brokers  []string `yaml:"brokers"`
}

// Service status config.
type Status struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	VersionPath   string `yaml:"versionPath"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
}

// Config contains all configuration parameters in the config package.
type Config struct {
	Project  Project  `yaml:"project"`
	GRPC     GRPC     `yaml:"grpc"`
	REST     REST     `yaml:"rest"`
	Database Database `yaml:"database"`
	Metrics  Metrics  `yaml:"metrics"`
	Jaeger   Jaeger   `yaml:"jaeger"`
	Kafka    Kafka    `yaml:"kafka"`
	Status   Status   `yaml:"status"`
}

// ReadConfigYML reads configurations from file and inits instance Config.
func ReadConfigYML(configYML string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(configYML)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return nil
}
