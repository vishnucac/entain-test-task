# Entain Test task

The main goal of this test task is a develop the application for processing the incoming requests from the 3d-party providers.
Technologies: Golang + Postgres.

## Requirements

- Docker
- Docker Compose

## Getting Started

### 1. Clone the repository
git clone https://github.com/vishnucac/entain-test-task.git
cd entain-test-task


### 2. Build and run the app with Docker Compose
docker-compose up -d

This command will start the application and the PostgreSQL database. The app will be available at `http://localhost:8080`.

### 3. Stopping the application
docker-compose down

## API Endpoints

- GET /status: Check if the service is running.
  - Response: `{"status": "service is up"}`

- GET /users: Get a list of users and their balances.
  - Response:
  [
    {"userId": 1, "balance": 100},
    {"userId": 2, "balance": 50.5},
    {"userId": 3, "balance": 200.75}
  ]

- POST /user/{userId}/transaction: Update a user’s balance based on the transaction state (`win` or `lose`).
  - Request body:
  {
    "state": "win", 
    "amount": 50.00,
    "transactionId": "1234"
  }

  - Response: `HTTP 200 OK`

- GET /user/{userId}/balance: Get a user’s balance.
  - Response:
  {
    "userId": 1,
    "balance": 100.00
  }

## Running Tests
go test ./..