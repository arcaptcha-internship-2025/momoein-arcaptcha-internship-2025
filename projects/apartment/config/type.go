package config

type Config struct {
	AppMode AppModeType  `json:"appMode" env:"APP_MODE"`
	DB      DBConfig     `json:"db"`
	HTTP    HTTPConfig   `json:"http"`
	Auth    AuthConfig   `json:"auth"`
	SMTP    SMTPConfig   `json:"smtp"`
	Minio   MinioConfig  `json:"minio"`
	BaseURL string       `json:"baseURL" env:"BASE_URL"`
	Smaila  SmailaConfig `json:"smaila"`
}

type AppModeType string

const (
	Development AppModeType = "development"
	Production  AppModeType = "production"
)

type DBConfig struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     uint   `json:"port" env:"DB_PORT"`
	DBName   string `json:"dbName" env:"DB_NAME"`
	Schema   string `json:"schema" env:"DB_SCHEMA"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASSWORD"`
	AppName  string `json:"appName" env:"DB_APP_NAME"`
}

type HTTPConfig struct {
	Port uint `json:"port" env:"HTTP_PORT"`
}

type AuthConfig struct {
	JWTSecret     string `json:"jwtSecret" env:"AUTH_JWT_SECRET"`
	AccessExpiry  int64  `json:"accessExpiry" env:"AUTH_ACCESS_EXPIRY"`
	RefreshExpiry int64  `json:"refreshExpiry" env:"AUTH_REFRESH_EXPIRY"`
}

type SMTPConfig struct {
	Host     string `json:"host" env:"SMTP_HOST"`
	Port     uint   `json:"port" env:"SMTP_PORT"`
	From     string `json:"from" env:"SMTP_FROM"`
	Username string `json:"username" env:"SMTP_USERNAME"`
	Password string `json:"password" env:"SMTP_PASSWORD"`
}

type MinioConfig struct {
	Endpoint  string `json:"endpoint" env:"MINIO_ENDPOINT"`
	AccessKey string `json:"accessKey" env:"MINIO_ACCESS_KEY"`
	SecretKey string `json:"secretKey" env:"MINIO_SECRET_KEY"`
}

type SmailaConfig struct {
	Endpoint string `json:"endpoint" env:"SMAILA_ENDPOINT"`
}
