markdown
# Go JSON Client-Server User Manager

A lightweight client-server application written in Go that demonstrates HTTP communication using JSON payloads and PostgreSQL database persistence.

## Overview

The project consists of two independent Go modules:
1. **`gojson` (Server):** An HTTP REST-like server built with Go's standard `net/http` package and `pgx/v5` connection pooling to store and retrieve user records in PostgreSQL.
2. **`clientjson` (CLI Client):** An interactive command-line client that communicates with the server via `GET` and `POST` requests to list existing users or register new ones.

---

## Project Structure

```text
.
├── clientjson/
│   ├── go.mod
│   └── main.go        # Interactive CLI client
└── gojson/
    ├── go.mod
    ├── go.sum
    └── main.go        # HTTP server and PostgreSQL integration
```

---

## Tech Stack

* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Database Driver:** [jackc/pgx/v5](https://github.com) (`pgxpool`)
* **Standard Library:** `net/http`, `encoding/json`, `bufio`, `context`

---

## Prerequisites

* **Go** 1.20 or higher
* **PostgreSQL** instance running and accessible

---

## Database Setup

1. Create a PostgreSQL database named `testjson` and run the following SQL query to create the `users` table:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone BIGINT NOT NULL
);
```

2. Open `gojson/main.go` and replace the placeholder credentials in `connStr` with your actual PostgreSQL connection details:

```go
connStr := "postgres://login:password@ip:port/testjson"
```

---

## Getting Started

### 1. Start the Server

Navigate to the server directory, download the required dependencies, and start the HTTP server:

```bash
cd gojson
go mod download
go run main.go
```

The server will start listening on `http://localhost:8888`.

### 2. Run the Client

Open a new terminal window, navigate to the client directory, and run the CLI tool:

```bash
cd clientjson
go run main.go
```

Follow the on-screen prompts:
* Enter **`1`** to fetch and display all registered users.
* Enter **`2`** to register a new user (name and numeric phone number).

---

## API Endpoints

| Method | Endpoint | Description | Request Body | Response |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/userlist` | Retrieves all users from the database | None | JSON array of `User` objects |
| `POST` | `/useradd` | Adds a new user to the database | `{"name": "John", "phone": 123456789}` | `201 Created` confirmation text |

### Example JSON Payload

```json
{
  "id": 1,
  "name": "Alice",
  "phone": 79991234567
}
```