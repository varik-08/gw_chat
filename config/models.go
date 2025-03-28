package config

type HTTP struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"`
	Port string `env:"PORT" envDefault:"8080"`
}

type DB struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	Schema   string `env:"SCHEMA"`
}

type Cfg struct {
	AppEnv    string `env:"ENV" envDefault:"development"`
	HTTP      HTTP   `envPrefix:"HTTP_"`
	DB        DB     `envPrefix:"DB_"`
	SecretKey string `env:"SECRET_KEY"`
}
