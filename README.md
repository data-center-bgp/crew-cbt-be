# Crew CBT (Competency Based Training) Web App - API Backend

A Go-based REST API backend for crew CBT, used for internal use only.

## Tech Stacks

- Go 1.22.1
- Fiber v2 (Web Framework)
- PostgreSQL (Database)
- GORM (ORM/Object Relational Mapper, for handling database)

## Project Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Environment variables configured in `.env` file

## Environment Variables

```env
DB_HOST=your_database_host
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=your_database_name
DB_PORT=your_database_port
DB_SSLMODE=disable
```

# API Endpoints

## User Management

- `POST /api/crew_cbt/user/register` - Register new crew for taking CBT
- `GET /api/crew_cbt/user/:nik` - Get crew quiz results based by NIK

## CBT Management

- `GET /api/crew_cbt/quiz/categories` - Get all quiz categories
- `GET /api/crew_cbt/quiz/categories/:id` - Get questions for specific category
- `POST /api/crew_cbt/quiz/submit` - Submit quiz answers
- `POST /api/crew_cbt/quiz/createCategory` - Create new quiz category
- `POST /api/crew_cbt/quiz/createQuestion` - Create new quiz question

# Running the application

1. Install dependencies
   `go mod download`

2. Start the server
   `go run main.go`

The server will start on port 3000 by default.

# CORS (Cross-Origin Resource Sharing) Configuration

The API is configured to allow requests from:

- `http://localhost:5173` (Development)
- `https://images.barokahperkasagroup.id` (Production)

# Features

- Crew registration and management for CBT
- Quiz categories management
- Question and answer management with image support and multiple answer management
- Quiz attempt tracking
- Score calculation and /pass/fail determination
- Time tracking for quiz attempts
