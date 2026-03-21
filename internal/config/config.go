package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ConfigHTTPServer struct {
	Host string `env:"HOST" env-default:"localhost"`
	Port string `env:"PORT" env-default:"8080"`
}

type ConfigDatabase struct {
	Role string `env:"DB_ROLE"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
}

type ConfigRedis struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Username string `env:"REDIS_USERNAME"`
	Pass     string `env:"REDIS_PASS"`
	DB       int    `env:"REDIS_DB"`
}

type ConfigRabbitMQ struct {
	Host     string `env:"RABBIT_HOST" default:"localhost"`
	Port     string `env:"RABBIT_PORT" default:"5672"`
	User     string `env:"RABBIT_USER" default:"guest"`
	Password string `env:"RABBIT_PASSWORD" default:"guest"`
	VHost    string `env:"RABBIT_VHOST" default:"/"`

	ExchangeName string `env:"RABBIT_EXCHANGE_NAME" default:"config.updates"`
	ExchangeType string `env:"RABBIT_EXCHANGE_TYPE" default:"fanout"`
	QueuePrefix  string `env:"RABBIT_QUEUE_PREFIX" default:"config-service."`
}

type ConfigLogger struct {
	APIlp   string `env:"APILOGFILENAME"`
	DBlp    string `env:"DBLOGFILENAME"`
	LogsDir string `env:"LOGSDIR"`
}

type Config struct {
	HTTPServer *ConfigHTTPServer
	Db         *ConfigDatabase
	Redis      *ConfigRedis
	Rabbit     *ConfigRabbitMQ
	Logger     *ConfigLogger
}

func MustLoad() *Config {
	godotenv.Load("./configs/.env")

	var cfgHttpServer ConfigHTTPServer
	var cfgDatabase ConfigDatabase
	var cfgRedis ConfigRedis
	var cfgRabbit ConfigRabbitMQ
	var cfgLogger ConfigLogger

	if err := cleanenv.ReadEnv(&cfgHttpServer); err != nil {
		log.Fatalln("Ошибка чтения настроек server из .env файлы")
	}
	if err := cleanenv.ReadEnv(&cfgDatabase); err != nil {
		log.Fatalln("Ошибка чтения настроек db из .env файлы")
	}
	if err := cleanenv.ReadEnv(&cfgRedis); err != nil {
		log.Fatalln("Ошибка чтения настроек redis из .env файла")
	}
	if err := cleanenv.ReadEnv(&cfgRabbit); err != nil {
		log.Fatalln("Ошибка чтения настроек rabbitmq из .env файла")
	}
	if err := cleanenv.ReadEnv(&cfgLogger); err != nil {
		log.Fatalln("Ошибка чтения настроек logger из .env файла")
	}

	return &Config{&cfgHttpServer, &cfgDatabase, &cfgRedis, &cfgRabbit, &cfgLogger}
}
