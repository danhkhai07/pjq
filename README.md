# pjq - Pluggable Job Queue

**pjq** is a robust, asynchronous job queue and worker pool service built in Go. It is designed to be easily extensible with a pluggable handler system and uses PostgreSQL for reliable, persistent storage of jobs.

## Features

- **REST API:** Simple HTTP endpoints to submit, retrieve, and filter jobs.
- **Persistent Storage:** PostgreSQL backing ensures jobs are not lost if the application restarts.
- **Pluggable Architecture:** Easily register new job handlers for different job types.
- **Worker Pool:** Processes jobs concurrently with an in-memory queue and managed worker routines.
- **Retry Mechanism:** Built-in support for job retries and tracking.
- **Status Tracking:** Comprehensive job states (`pending`, `running`, `done`, `failed`, `retrying`).

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Docker & Docker Compose (optional, for easy setup)

### Running Locally

1. **Start the database:**
   If you have Docker Compose installed, you can start the PostgreSQL database easily:
   ```bash
   docker-compose up -d
   ```

2. **Run the application:**
   ```bash
   go run cmd/main.go
   ```
   The application will start on `http://0.0.0.0:8888`.

## API Documentation

### 1. Welcome / Health Check
Check if the service is running.

- **URL:** `GET /`
- **Response:**
  ```json
  {
    "message": "Welcome to pjq!"
  }
  ```

### 2. Submit a Job
Enqueues a new job for background processing.

- **URL:** `POST /jobs`
- **Body:**
  ```json
  {
    "type": "email_send",
    "payload": {
      "to": "user@example.com",
      "subject": "Hello!"
    }
  }
  ```
- **Response:** (201 Created)
  ```json
  {
    "id": "01H...XYZ"
  }
  ```

### 3. Get Job Details
Retrieve the status, payload, and result of a specific job.

- **URL:** `GET /jobs/{id}`
- **Response:**
  ```json
  {
    "id": "01H...XYZ",
    "type": "email_send",
    "status": "done",
    "result": "Email sent successfully",
    "created_at": "2026-06-01T12:00:00Z"
  }
  ```

### 4. List / Filter Jobs
List all jobs, optionally filtering by status, type, or whether they are retriable.

- **URL:** `GET /jobs`
- **Body:** (Optional, as JSON)
  ```json
  {
    "status": "pending",
    "type": "email_send"
  }
  ```
- **Response:**
  ```json
  {
    "total": 1,
    "jobs": [
      {
        "id": "01H...XYZ",
        "type": "email_send",
        "status": "pending"
      }
    ]
  }
  ```

## Architecture

- **Domain:** Defines core entities like `Job` and states.
- **Queue & Registry:** Manages the in-memory queue, concurrency, and dispatching jobs to registered handlers based on the job `Type`.
- **Infrastructure:** Implements the PostgreSQL store and concrete handlers.
- **API & Application:** Provides REST endpoints and coordinates between the API requests and the underlying queue/storage mechanisms.
