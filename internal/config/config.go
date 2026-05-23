package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ConfigHTTPServer struct {
	Host            string        `env:"HTTPHOST" env-default:"localhost"`
	Port            int           `env:"HTTPPORT" env-default:"8080"`
	ShutdownTimeout time.Duration `env:"HTTPSHUTDOWNTIMEOUT" env-default:"30s"`
}

type ConfigGRPCServer struct {
	Host            string        `env:"GRPCHOST" env-default:"localhost"`
	Port            int           `env:"GRPCPORT" env-default:"8080"`
	ShutdownTimeout time.Duration `env:"GRPCSHUTDOWNTIMEOUT" env-default:"30s"`
}

type ConfigDatabase struct {
	Role string `env:"DB_ROLE"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
}

type ConfigLogger struct {
	APIlp   string `env:"APILOGFILENAME"`
	DBlp    string `env:"DBLOGFILENAME"`
	GRPClp  string `env:"GRPCLOGFILENAME"`
	LogsDir string `env:"LOGSDIR"`
}

type ConfigSchema struct {
	SchemasDir string `env:"SCHEMASDIR"`
}

type Config struct {
	HTTPServer *ConfigHTTPServer
	GRPCServer *ConfigGRPCServer
	Db         *ConfigDatabase
	Logger     *ConfigLogger
	Schema     *ConfigSchema
}

func MustLoad() *Config {
	if err := godotenv.Load("./configs/.env"); err != nil {
		log.Fatalf("file has not found: %v", err)
	}

	var cfgHTTPServer ConfigHTTPServer
	var cfgGRPCServer ConfigGRPCServer
	var cfgDatabase ConfigDatabase
	var cfgLogger ConfigLogger
	var cfgSchema ConfigSchema

	if err := cleanenv.ReadEnv(&cfgHTTPServer); err != nil {
		log.Fatalln("Ошибка чтения настроек server из .env файлы")
	}
	if err := cleanenv.ReadEnv(&cfgGRPCServer); err != nil {
		log.Fatalln("Ошибка чтения настроек grpc из .env файлы")
	}
	if err := cleanenv.ReadEnv(&cfgDatabase); err != nil {
		log.Fatalln("Ошибка чтения настроек db из .env файлы")
	}
	if err := cleanenv.ReadEnv(&cfgLogger); err != nil {
		log.Fatalln("Ошибка чтения настроек logger из .env файла")
	}
	if err := cleanenv.ReadEnv(&cfgSchema); err != nil {
		log.Fatalln("Ошибка чтения настроек schema из .env файла")
	}

	return &Config{&cfgHTTPServer, &cfgGRPCServer, &cfgDatabase, &cfgLogger, &cfgSchema}
}
