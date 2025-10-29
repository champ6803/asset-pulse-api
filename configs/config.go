package config

type Config struct {
	AppEnv              string `env:"APP_ENV"`
	AppPort             string `env:"APP_PORT"`
	AppName             string `env:"APP_NAME"`
	HttpTimeOut         int    `env:"HTTP_TIMEOUT"`
	HttpRetryCount      int    `env:"HTTP_RETRY_COUNT"`
	JWTSecret           string `env:"JWT_SECRET"`
	DbHost              string `env:"DB_HOST"`
	DbPort              string `env:"DB_PORT"`
	DbName              string `env:"DB_NAME"`
	DbUsername          string `env:"DB_USERNAME"`
	DbPassword          string `env:"DB_PASSWORD"`
	DbSSLMode           string `env:"SSL_MODE"`
	DbSSLRootCert       string `env:"SSL_ROOT_CERT"`
	DbTimeZone          string `env:"TIME_ZONE"`
	AzureOpenAIKey      string `env:"AZURE_OPENAI_KEY"`
	AzureOpenAIEndpoint string `env:"AZURE_OPENAI_ENDPOINT"`
	AzureOpenAIModel    string `env:"AZURE_OPENAI_MODEL" envDefault:"gpt-4"`
	OpenAIAPIKey        string `env:"OPENAI_API_KEY"`
	// remove later
	AzureOpenaiModelName string `env:"AZURE_OPENAI_MODEL_NAME"`
	AzureOpenaiUrl       string `env:"AZURE_OPENAI_URL"`
}
