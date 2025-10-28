package config

type Config struct {
	AppEnv         string `env:"APP_ENV"`
	AppPort        string `env:"APP_PORT"`
	AppName        string `env:"APP_NAME"`
	HttpTimeOut    int    `env:"HTTP_TIMEOUT"`
	HttpRetryCount int    `env:"HTTP_RETRY_COUNT"`
	JWTSecret      string `env:"JWT_SECRET"`
	OpenAIAPIKey   string `env:"OPENAI_API_KEY"`
	DbHost         string `env:"DB_HOST"`
	DbPort         string `env:"DB_PORT"`
	DbName         string `env:"DB_NAME"`
	DbUsername     string `env:"DB_USERNAME"`
	DbPassword     string `env:"DB_PASSWORD"`
	DbSSLMode      string `env:"SSL_MODE"`
	DbSSLRootCert  string `env:"SSL_ROOT_CERT"`
	DbTimeZone     string `env:"TIME_ZONE"`
}
