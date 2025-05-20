# 🎟️ Go Asynq Ticket System Example

A simple asynchronous ticket queue system built with Go, Echo (HTTP framework), and Asynq (Redis-based task queue).  
This project illustrates how to enqueue tasks via an HTTP API and process them in the background.

## Table of Contents
- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Clone the Repository](#clone-the-repository)
  - [Configure Redis](#configure-redis)
  - [Initialize Go Modules](#initialize-go-modules)
- [Running the Services](#running-the-services)
  - [Start the Echo API](#start-the-echo-api)
  - [Start the Echo Worker](#start-the-echo-worker)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Configuration](#configuration)
- [Docker (Optional)](#docker-optional)
- [Contributing](#contributing)
- [License](#license)

## Overview

The system is composed of two components:
- **echo-api**: An HTTP server to accept ticket purchase requests and enqueue them as Asynq tasks.
- **echo-worker**: A background worker that processes queued tasks and simulates ticket sales.

## Prerequisites

- Go 1.21 or newer
- Redis server (default address: `127.0.0.1:6370`)
- (Optional) Docker & Docker Compose

## Project Structure

```
go-asynq-test/
├── echo-api/      # HTTP API service for creating tasks
└── echo-worker/   # Background worker service to process tasks
```

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/ShognDev2409/go-asynq-test.git
cd go-asynq-test
```

### Configure Redis

Make sure a Redis instance is running at `127.0.0.1:6370`. To start one using Docker:

```bash
docker run --name go-asynq-redis -p 6370:6379 -d redis:7
```

### Initialize Go Modules

```bash
cd echo-api
go mod tidy

cd ../echo-worker
go mod tidy
```

## Running the Services

### Start the Echo API

In a terminal:
```bash
cd echo-api
go run main.go --http ":8080" --redis "127.0.0.1:6370"
```
The API will listen on port `8080` by default.

### Start the Echo Worker

In another terminal:
```bash
cd echo-worker
go run main.go --redis "127.0.0.1:6370" --concurrency 10
```
The worker will poll the queue and process tasks concurrently (default 10 workers).

## API Endpoints

### POST /tickets

Enqueue a ticket purchase task.

**Request Body** (JSON):
```json
{
  "user_id": "string",
  "event_id": "string",
  "quantity": 1
}
```

**Response** (JSON):
```json
{
  "status": "queued",
  "task_id": "string"
}
```

## Testing

Run unit and concurrency tests in the `echo-worker` module:

```bash
cd echo-worker
go test ./... -v
```

## Configuration

| Environment   | Flag           | Default            | Description                     |
|--------------:|:---------------|:-------------------|:--------------------------------|
| `API_PORT`     | `--http`       | `:8080`            | HTTP server listen address      |
| `REDIS_ADDR`   | `--redis`      | `127.0.0.1:6370`   | Redis server address            |
| `CONCURRENCY`  | `--concurrency`| `10`               | Max concurrent worker handlers  |

## Docker (Optional)

A Docker Compose setup can simplify running Redis, API, and Worker together. Example `docker-compose.yml` not included.

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues](https://github.com/ShognDev2409/go-asynq-test/issues) and submit pull requests.

## License

This project is licensed under the MIT License.
