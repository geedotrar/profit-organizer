# Laba Service (Go + Gin + GORM)
## ✨ Features

- Get All Laba Data
- Import Data Laba From Excel File
- Export Data Laba to Excel File

## 🧱 Tech Stack
- Docker & Docker Compose
- Language : Go
- Framework : Gin
- ORM : GORM
- Database : PostgreSQL
- Env Config : godotenv

## 🚀 Setup
1. **Navigate to the Laba Service directory**
   ```
   cd backend
   ```
2. **Copy the example environment file**
   ```
   cp .env.example .env
   ```
3. **Run the database**
   ```
   docker-compose up -d
   ```
4. **Install Go dependencies**
   ```
   go mod tidy
   ```
5. **Run database migrations and seeders**
   ```
   make migrate-up
   ```
6. **Start the service**
   ```
   go run cmd/main.go
   ```

## 🛣️ API Endpoints

| Method | Endpoint         | Description               |
|--------|------------------|---------------------------|
| GET    | `/laba`          | Get all laba records      |
| POST   | `/laba/import`   | Import laba from Excel    |
| GET    | `/laba/export`   | Export laba to Excel file |
