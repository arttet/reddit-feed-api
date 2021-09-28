package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
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

// Project contains all parameters project information.
type Project struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	SwaggerDir  string `yaml:"swaggerDir"`
	Version     string
	CommitHash  string
}

// Database contains all parameters database connection.
type Database struct {
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Name          string `yaml:"name"`
	SslMode       string `yaml:"sslmode"`
	MigrationsDir string `yaml:"migrationsDir"`
	Driver        string `yaml:"driver"`
}

// GRPC contains all parameters of gRPC.
type GRPC struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	MaxConnectionIdle int64  `yaml:"maxConnectionIdle"`
	Timeout           int64  `yaml:"timeout"`
	MaxConnectionAge  int64  `yaml:"maxConnectionAge"`
}

// REST contains all parameters of REST.
type REST struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Jaeger contains all parameters tracer information.
type Jaeger struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Service string `yaml:"service"`
}

// Metrics contains all parameters metrics information.
type Metrics struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// Status contains all parameters status information.
type Status struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	LivenessPath  string `yaml:"livenessPath"`
	ReadinessPath string `yaml:"readinessPath"`
	SwaggerPath   string `yaml:"swaggerPath"`
	VersionPath   string `yaml:"versionPath"`
}

// Kafka contains all parameters Kafka information.
type Kafka struct {
	Capacity int      `yaml:"capacity"`
	Topic    string   `yaml:"topic"`
	GroupID  string   `yaml:"groupId"`
	Brokers  []string `yaml:"brokers"`
}

// Config contains all configuration parameters in the config package.
type Config struct {
	Project  Project    `yaml:"project"`
	Database Database   `yaml:"database"`
	Logger   zap.Config `yaml:"logger"`
	GRPC     GRPC       `yaml:"grpc"`
	REST     REST       `yaml:"rest"`
	Jaeger   Jaeger     `yaml:"jaeger"`
	Metrics  Metrics    `yaml:"metrics"`
	Status   Status     `yaml:"status"`
	Kafka    Kafka      `yaml:"kafka"`
}

// String returns a data source name.
func (db *Database) String() string {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		db.Host,
		db.Port,
		db.User,
		db.Password,
		db.Name,
		db.SslMode,
	)

	return dsn
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
