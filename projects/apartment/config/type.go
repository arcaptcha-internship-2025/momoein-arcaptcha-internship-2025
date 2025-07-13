package config

type Config struct {
	AppMode AppModeType `json:"appMode"`
	DB      DBConfig    `json:"db"`
	HTTP    HTTPConfig  `json:"http"`
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
