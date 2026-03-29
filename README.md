# Cinema Booking System

A concurrent seat booking system built in Go with Redis backend. The system efficiently handles concurrent seat reservations for multiple movies with automatic session expiration.

## Features

- **Concurrent Booking**: Thread-safe seat reservation handling multiple simultaneous requests
- **Session-based Management**: Temporary hold periods for seats before confirmation
- **Redis-backed Storage**: Distributed session management with automatic TTL expiration
- **RESTful API**: Clean HTTP interface for all booking operations
- **Web UI**: Static HTML interface for browsing movies and seats

## Architecture

### Key Components

- **Service Layer**: Business logic for booking operations
- **Store Interfaces**: Multiple backend implementations (Redis, Memory, Concurrent)
  - `RedisStore`: Production-ready distributed session store
  - `MemoryStore`: In-memory implementation for testing
  - `ConcurrentStore`: Thread-safe in-memory store with RWMutex
- **HTTP Handlers**: RESTful API endpoints
- **Redis Adapter**: Connection management and client initialization

### Booking Flow

1. **Hold Seat**: User initiates a booking → Seat is reserved with 2-minute TTL
2. **Confirm Session**: User confirms the booking → TTL is removed, booking becomes permanent
3. **Release Session**: User cancels the booking → Seat is released back to available pool

## Prerequisites

- **Go**: 1.26.1 or later
- **Docker & Docker Compose**: For running Redis
- **Redis**: 7-alpine (managed via Docker Compose)

## Installation

### 1. Clone/Navigate to Project

```bash
cd "Cinema Booking System GoLang"
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start Redis

```bash
docker compose up -d
```

Verify Redis is running:
```bash
docker compose ps
```

## Running the Application

### Development Mode

```bash
go run ./cmd/main.go
```

### Build and Run

```bash
go build -o cinema ./cmd/main.go
./cinema
```

The application will start on **http://localhost:8080**

### Expected Output

```
2026/03/29 19:47:24 connected to redis at localhost:6565
```

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── adapters/
│   │   └── redis/
│   │       └── redis.go        # Redis client initialization
│   └── booking/
│       ├── domain.go           # Core domain models (Booking, BookingStore interface)
│       ├── service.go          # Business logic and orchestration
│       ├── handler.go          # HTTP request handlers
│       ├── memory_store.go     # In-memory booking store
│       ├── redis_store.go      # Redis-backed booking store
│       ├── concurrent_store.go # Thread-safe in-memory store
│       └── service_test.go     # Service tests
├── static/
│   └── index.html              # Web UI
├── docker-compose.yaml         # Docker Compose configuration
├── go.mod                       # Go module dependencies
└── README.md                    # This file
```

## API Endpoints

### Movies

#### List Available Movies
```
GET /movies
```

**Response:**
```json
[
  {
    "id": "inception",
    "title": "Inception",
    "rows": 5,
    "seats_per_row": 8
  },
  {
    "id": "dune",
    "title": "Dune: Part Two",
    "rows": 4,
    "seats_per_row": 6
  }
]
```

### Seat Management

#### List Seats for a Movie
```
GET /movies/{movieID}/seats
```

**Response:**
```json
[
  {
    "seat_id": "A1",
    "user_id": "user123",
    "booked": true,
    "confirmed": false
  }
]
```

#### Hold a Seat (Create Session)
```
POST /movies/{movieID}/seats/{seatID}/hold
Content-Type: application/json

{
  "user_id": "user123"
}
```

**Response:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "movieID": "inception",
  "seat_id": "A1",
  "expires_at": "2026-03-29T19:50:24Z"
}
```

**Status**: `201 Created`

### Session Management

#### Confirm a Session
```
PUT /sessions/{sessionID}/confirm
Content-Type: application/json

{
  "user_id": "user123"
}
```

**Response:**
```json
{
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "movie_id": "inception",
  "seat_id": "A1",
  "user_id": "user123",
  "status": "confirmed"
}
```

**Status**: `200 OK`

#### Release a Session
```
DELETE /sessions/{sessionID}
Content-Type: application/json

{
  "user_id": "user123"
}
```

**Status**: `204 No Content`

## Configuration

### Docker Compose

The application connects to Redis via Docker Compose. Port mapping:
- **Host**: `6565`
- **Container**: `6379` (Redis default)

Connection string: `localhost:6565`

### Default Settings

- **Server Port**: 8080
- **Hold TTL**: 2 minutes (120 seconds)
- **Movies**: Inception, Dune: Part Two

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

The `service_test.go` includes tests for concurrent booking scenarios.

## Technologies Used

- **Language**: Go 1.26.1
- **Database**: Redis 7.0
- **HTTP Server**: Go net/http (standard library)
- **Dependencies**:
  - `github.com/redis/go-redis/v9` - Redis client
  - `github.com/google/uuid` - UUID generation
- **Containerization**: Docker & Docker Compose

## Development Guide

### Adding a New Store Implementation

Implement the `BookingStore` interface in `internal/booking/domain.go`:

```go
type BookingStore interface {
    Book(b Booking) (Booking, error)
    ListBookings(movieID string) []Booking
    Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
    Release(ctx context.Context, sessionID string, userID string) error
}
```

### Key Constants

- `defaultHoldTTL`: Session hold duration before automatic expiration (2 minutes)
- `Movies`: Hard-coded in `cmd/main.go`

## Troubleshooting

### Redis Connection Error

If you see `redis ping: EOF`:

1. Ensure Docker is running
2. Restart Redis:
   ```bash
   docker compose restart redis
   ```
3. Check port mapping is correct (6565 → 6379)

### Port Already in Use

If port 8080 is occupied:

1. Find the process: `lsof -i :8080`
2. Kill it or modify the port in `cmd/main.go` line 31

### Build Errors

Ensure all dependencies are installed:

```bash
go mod tidy
go mod download
```

## Future Enhancements

- [ ] Database persistence (PostgreSQL)
- [ ] User authentication & authorization
- [ ] Payment processing
- [ ] Email notifications
- [ ] Admin dashboard
- [ ] Analytics and reporting
- [ ] Seat selection UI with visual layout
- [ ] Multi-cinema support

## License

MIT

## Author

Cinema Booking System Team