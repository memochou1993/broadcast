distributed-chatroom
===

# Usage

Run a Redis database.

```BASH
docker run --name redis -d -p 6379:6379 redis --requirepass password
```

Copy `.env.example` to `.env`.

```ENV
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=password
```

Run the server.

```BASH
go run ./
```

Or run the server with docker compose.

```BASH
docker-compose up -d
```

Visit <http://localhost>
