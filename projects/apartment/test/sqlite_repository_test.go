package test

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/arcaptcha-internship-2025/momoein-apartment/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	log = logger.NewConsoleZapLogger(logger.ModeDevelopment)
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	wd, _ := os.Getwd()
	log.Debug("", zap.Any("wd", wd))

	schemaBytes, err := os.ReadFile("../sql/schema.sqlite.sql")
	require.NoError(t, err)

	schema := string(schemaBytes)
	for _, stmt := range strings.Split(schema, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			_, err := db.Exec(stmt)
			require.NoError(t, err, "failed executing SQL: %s", stmt)
		}
	}

	return db
}

func TestSetupTestDB(t *testing.T) {
	db := setupTestDB(t)
	require.NotNil(t, db)
	defer db.Close()

	err := db.QueryRow(`SELECT * FROM users;`).Err()
	assert.Nil(t, err)
}
