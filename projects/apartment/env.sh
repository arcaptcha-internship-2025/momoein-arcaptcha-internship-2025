#!/bin/bash

set -e

MODE=${1:-interactive}

# Default values

APP_MODE="development"
BASE_URL="http://127.0.0.1:8080"
HTTP_PORT="8080"

DB_HOST="apartment-db"
DB_PORT="5432"
DB_NAME="postgres"
DB_SCHEMA="public"
DB_USER="postgres"
DB_PASSWORD="postgres"
DB_APP_NAME="apartment-api"

AUTH_JWT_SECRET="I am the secret skyler"
AUTH_ACCESS_EXPIRY="1440"
AUTH_REFRESH_EXPIRY="14400"

MINIO_ENDPOINT="apartment-minio:9000"
MINIO_ACCESS_KEY="minioadmin"
MINIO_SECRET_KEY="minioadmin"

SMAILA_ENDPOINT="http://apartment-smaila:1174"

# smaila config
SMAILA_HTTP_PORT="1174"
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_FROM="apartment.arcaptcha@gmail.com"
SMTP_USERNAME="apartment.arcaptcha@gmail.com"
SMTP_PASSWORD=""

function prompt_with_default() {
    local prompt_text=$1
    local default_val=$2
    local varname=$3
    local input

    read -p "$prompt_text [$default_val]: " input
    input=${input:-$default_val}
    printf -v "$varname" '%s' "$input"
}

function prompt_password() {
    local prompt_text=$1
    local varname=$2
    local input=""
    while [[ -z "$input" ]]; do
        read -sp "$prompt_text: " input
        echo ""
    done
    printf -v "$varname" '%s' "$input"
}

echo "Running setup in $MODE mode..."

if [[ "$MODE" == "interactive" ]]; then

    # App mode with validation
    while true; do
        echo "Choose app mode: (development / production)"
        read -p "App mode [${APP_MODE}]: " input
        input=${input:-$APP_MODE}
        if [[ "$input" == "development" || "$input" == "production" ]]; then
            APP_MODE=$input
            break
        else
            echo "❌ Invalid choice. Please enter 'development' or 'production'."
        fi
    done

    prompt_with_default "Base URL" "$BASE_URL" BASE_URL
    prompt_with_default "HTTP port" "$HTTP_PORT" HTTP_PORT

    # DB config
    prompt_with_default "DB host" "$DB_HOST" DB_HOST
    prompt_with_default "DB port" "$DB_PORT" DB_PORT
    prompt_with_default "DB name" "$DB_NAME" DB_NAME
    prompt_with_default "DB schema" "$DB_SCHEMA" DB_SCHEMA
    prompt_with_default "DB user" "$DB_USER" DB_USER
    prompt_password "DB password (required)" DB_PASSWORD
    prompt_with_default "DB app name" "$DB_APP_NAME" DB_APP_NAME

    # Auth config
    prompt_password "JWT secret (required)" AUTH_JWT_SECRET
    prompt_with_default "Access expiry minutes" "$AUTH_ACCESS_EXPIRY" AUTH_ACCESS_EXPIRY
    prompt_with_default "Refresh expiry minutes" "$AUTH_REFRESH_EXPIRY" AUTH_REFRESH_EXPIRY

    # Minio config
    prompt_with_default "Minio endpoint" "$MINIO_ENDPOINT" MINIO_ENDPOINT
    prompt_with_default "Minio access key" "$MINIO_ACCESS_KEY" MINIO_ACCESS_KEY
    prompt_with_default "Minio secret key" "$MINIO_SECRET_KEY" MINIO_SECRET_KEY

    # Smaila config
    prompt_with_default "Smaila endpoint" "$SMAILA_ENDPOINT" SMAILA_ENDPOINT

    # SMTP config
    prompt_with_default "SMTP HTTP port" "$SMAILA_HTTP_PORT" SMAILA_HTTP_PORT
    prompt_with_default "SMTP host" "$SMTP_HOST" SMTP_HOST
    prompt_with_default "SMTP port" "$SMTP_PORT" SMTP_PORT
    prompt_with_default "SMTP from email" "$SMTP_FROM" SMTP_FROM
    prompt_with_default "SMTP username" "$SMTP_USERNAME" SMTP_USERNAME
    prompt_password "SMTP password (required)" SMTP_PASSWORD

elif [[ "$MODE" == "quick" ]]; then
    echo "Using default values except for required secrets."

    # Prompt only for required secrets
    prompt_password "SMTP password (required)" SMTP_PASSWORD
    # prompt_password "DB password (required)" DB_PASSWORD
    # prompt_password "JWT secret (required)" AUTH_JWT_SECRET

else
    echo "Unknown mode: $MODE"
    echo "Usage: $0 [interactive|quick]"
    exit 1
fi

# Write .smaila.env
cat > .smaila.env <<EOL
HTTP_PORT=${SMAILA_HTTP_PORT}
SMTP_HOST=${SMTP_HOST}
SMTP_PORT=${SMTP_PORT}
SMTP_FROM=${SMTP_FROM}
SMTP_USERNAME=${SMTP_USERNAME}
SMTP_PASSWORD=${SMTP_PASSWORD}
EOL
echo ".smaila.env file created."

# Write main .env
cat > .env <<EOL
APP_MODE=${APP_MODE}
BASE_URL=${BASE_URL}

# http config
HTTP_PORT=${HTTP_PORT}

# db config
DB_HOST=${DB_HOST}
DB_PORT=${DB_PORT}
DB_NAME=${DB_NAME}
DB_SCHEMA=${DB_SCHEMA}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_APP_NAME=${DB_APP_NAME}

# auth config
AUTH_JWT_SECRET=${AUTH_JWT_SECRET}
AUTH_ACCESS_EXPIRY=${AUTH_ACCESS_EXPIRY}
AUTH_REFRESH_EXPIRY=${AUTH_REFRESH_EXPIRY}

# minio config
MINIO_ENDPOINT=${MINIO_ENDPOINT}
MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY}
MINIO_SECRET_KEY=${MINIO_SECRET_KEY}

# smaila config
SMAILA_ENDPOINT=${SMAILA_ENDPOINT}
EOL

echo ".env file created successfully."
