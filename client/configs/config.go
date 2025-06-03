package configs

import (
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Api      Api
	Database DatabaseConfig
	RabbitMQ RabbitMQ
}

type Api struct {
	Port   int
	Queues map[string]string
}

type RabbitMQ struct {
	Host     string
	Port     int
	User     string
	Password string
	Queues   map[string]RabbitMQQueue
}

type RabbitMQQueue struct {
	Name        string
	Exchange    string
	RoutingKey  string
	Durable     bool
	AutoAck     bool
	Exclusive   bool
	AutoDelete  bool
	Passive     bool
	ConsumerTag string
	Args        map[string]interface{}
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SSLMode  bool
	Timeout  int
}

func (db *DatabaseConfig) GetConnectionString() string {
	// Formato aceito pelo Go: "host=... port=... user=... password=... dbname=... sslmode=... connect_timeout=..."
	var builder strings.Builder
	builder.WriteString("host=")
	builder.WriteString(db.Host)
	builder.WriteString(" port=")
	builder.WriteString(strconv.Itoa(db.Port))
	builder.WriteString(" user=")
	builder.WriteString(db.User)
	builder.WriteString(" password=")
	builder.WriteString(db.Password)
	builder.WriteString(" dbname=")
	builder.WriteString(db.DbName)
	builder.WriteString(" sslmode=")
	if db.SSLMode {
		builder.WriteString("require")
	} else {
		builder.WriteString("disable")
	}
	builder.WriteString(" connect_timeout=")
	if db.Timeout > 0 {
		builder.WriteString(strconv.Itoa(db.Timeout))
	} else {
		builder.WriteString("30")
	}
	return builder.String()
}

func (r *RabbitMQ) GetConnectionString() string {
	var builder strings.Builder
	builder.WriteString("amqp://")
	builder.WriteString(r.User)
	builder.WriteString(":")
	builder.WriteString(r.Password)
	builder.WriteString("@")
	builder.WriteString(r.Host)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(r.Port))
	return builder.String()
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("../../configs")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
