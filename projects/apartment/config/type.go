package config

type Config struct {
	AppMode AppModeType `json:"appMode"`
	DB      DBConfig    `json:"db"`
	HTTP    HTTPConfig  `json:"http"`
	Auth    AuthConfig  `json:"auth"`
	SMTP    SMTPConfig  `json:"smtp"`
	Minio   MinioConfig `json:"minio"`
}

type AppModeType string

const (
	Development AppModeType = "development"
	Production  AppModeType = "production"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	DBName   string `json:"dbName"`
	Schema   string `json:"schema"`
	User     string `json:"user"`
	Password string `json:"password"`
	AppName  string `json:"appName"`
}

type HTTPConfig struct {
	Port uint `json:"port"`
}

type AuthConfig struct {
	JWTSecret     string `json:"jwtSecret"`
	AccessExpiry  int64  `json:"accessExpiry"`
	RefreshExpiry int64  `json:"refreshExpiry"`
}

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	From     string `json:"from"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MinioConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}
