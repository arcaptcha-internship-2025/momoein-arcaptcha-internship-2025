package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DBConnOptions struct {
	User    string
	Pass    string
	Host    string
	Port    uint
	DBName  string
	Schema  string
	AppName string
}

func (o DBConnOptions) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s application_name=%s sslmode=disable",
		o.Host, o.Port, o.User, o.Pass, o.DBName, o.Schema, o.AppName)
}

func NewPSQLConn(opt DBConnOptions) (*sql.DB, error) {
	return sql.Open("postgres", opt.PostgresDSN())
}
