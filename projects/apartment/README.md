# Arcaptcha Apartment API

A Go-based apartment management API service featuring billing, user management, payment processing, and file storage integrations.

---

## 🚀 Features

- **Apartment Management** – Create, update, and manage apartments.
- **Billing System** – Generate and store bills with object storage support.
- **User Management** – Secure authentication and JWT-based authorization.
- **Payment Processing** – Mock payment gateway for testing and integration.
- **File Storage** – MinIO S3-compatible object storage integration.
- **Email Notifications** – Powered by Smaila SMTP service.
- **Interactive API Docs** – Swagger UI included.

---

## ⚙️ Setup

### 1️⃣ Environment Setup

Create your .env file in one of the following ways:

Interactive:

```bash
make env
```

Quick with defaults (only secrets required):

```bash
make env-quick
```

Manual:

```bash
cp example.env .env
```

> Adjust values as needed, especially secrets.

---

### 2️⃣ Start Services with Docker Compose

```bash
docker-compose up -d --build
```

This will start:

- **PostgreSQL** (Database)
- **Smaila** (SMTP service)
- **MinIO** (Object storage)
- **Apartment API** (Your app)

---

### 3️⃣ Apply Database Schema (First Run)

Production Schema:

```bash
make migrate-db
```

---

### 4️⃣ View API Logs

```bash
docker logs -f apartment-api
```

## 📘 API Documentation

Once the API is running, access Swagger UI at:

[http://127.0.0.1:8080/api/v1/docs/swagger](http://127.0.0.1:8080/api/v1/docs/swagger)

---

## 📤 Example Usage

Test the API with curl:

```bash
curl http://localhost:8080/
```

---

## 🛠 Development

Run DB migrations:

```bash
make migrate-db-dev
```

> ⚠️ This command drops all data in the database before applying migrations. Use only in development environments!

Regenerate Swagger docs:

```bash
make swagger
```

---

## 📝 Environment Variables

See `example.env` and `example.smaila.env` for full list.

Key variables:

| Variable        | Description                 | Default           |
| --------------- | --------------------------- | ----------------- |
| APP_MODE        | App mode (development/prod) | development       |
| HTTP_PORT       | API server port             | 8080              |
| AUTH_JWT_SECRET | JWT signing secret          | required          |
| DB\_\*          | Database config             | -                 |
| MINIO\_\*       | MinIO S3 config             | defaults provided |
| SMAILA\_\*      | Smaila SMTP config          | defaults provided |

---

## 📝 Todo

- Implement real payment gateway integration
- Add apartment search filters
- Enhance unit test coverage
- Add CI/CD pipeline

---
