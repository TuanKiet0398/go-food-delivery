# food-delivery

A Go learning project building a food delivery backend, following a course curriculum from UI to database through to gRPC.

## Requirements

- Go 1.26+
- MySQL

## Installation

```bash
go mod download
```

Create a `.env` file in the project root with the following content:

```
MYSQL_CONN_STRING=user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

## Usage

```bash
go run main.go
```

The server listens on `http://localhost:8080` by default.

## API

### Health check

- `GET /ping` — returns `{"message": "pong"}`

### Restaurants (`/v1/restaurants`)

| Method | Path                    | Description                              |
| ------ | ----------------------- | ----------------------------------------- |
| POST   | `/v1/restaurants`       | Create a new restaurant                   |
| GET    | `/v1/restaurants`       | List restaurants (paginated)              |
| GET    | `/v1/restaurants/:id`   | Get a restaurant by id                    |
| PATCH  | `/v1/restaurants/:id`   | Update a restaurant's `name` and/or `addr`|
| DELETE | `/v1/restaurants/:id`   | Delete a restaurant by id                 |

**Pagination query params** (for `GET /v1/restaurants`):

- `page` — page number, default `1`
- `limit` — items per page, default `50`

**Filter query params** (for `GET /v1/restaurants`):

- `owner_id` — restrict results to a single owner

List responses are wrapped in `common.NewSuccessResponse(data, paging, filter)`:

```json
{ "data": [...], "paging": { "page": 1, "limit": 50, "total": 0 }, "filter": { "owner_id": 0 } }
```

Other successful responses are wrapped in `common.SimpleSuccessResponse`:

```json
{ "data": ... }
```

Deleting a restaurant is a soft delete: rows have an integer `status` column
(`1` = active, `0` = soft-deleted), and deleting an already-deleted restaurant
returns an error instead of deleting again.

Restaurant `id`s exposed over the API are base58-encoded virtual UIDs
(`common.UID`), not raw integers: a UID packs `localID` (the real DB id),
`objectType` (the entity's DB type, e.g. `common.DbTypeRestaurant`) and
`shardID` into a single 64-bit value. `DELETE /v1/restaurants/:id` decodes
the UID back into the local DB id via `common.FromBase58`.

## Error handling

Handlers no longer write ad-hoc `c.JSON(http.StatusBadRequest, ...)` error
bodies. Business/storage errors are wrapped into a `*common.AppError`
(`common.ErrInvalidRequest`, `common.ErrDB`, `common.ErrCannotCreateEntity`,
etc.) and `panic`'d from the transport layer. `middleware.Recover`, mounted
in `main.go`, recovers the panic and writes the appropriate JSON status and
body — an `*common.AppError` keeps its own status code/message, any other
recovered value is wrapped as a `500` via `common.ErrInternal`.

## Architecture

Handlers no longer take a `*gorm.DB` directly. A shared `appctx.AppContext`
(`component/appctx`) wraps the DB connection and is threaded through
`main.go` into the transport layer, so future shared dependencies (config,
logger, etc.) can be added in one place.

See [architecture.md](architecture.md) for the full layered structure
(transport → biz → storage → model), directory layout, and request-flow
walkthroughs.

## Roadmap

- [x] Section 02 - UI to Database
- [x] Section 03 - GORM
- [x] Section 04 - Simple clean architecture
- [x] Section 05 - Error handling and UID
- [ ] Section 06 - Upload file to AWS S3
- [ ] Section 07 - Authenticate with JWT
- [ ] Section 08 - Linking model user and repository layer
- [ ] Section 09 - Like and dislike restaurant
- [ ] Section 10 - Async job, group and Pub/Sub
- [ ] Section 11 - Q&A mentoring, review CV
- [ ] Section 12 - Pub/Sub continue and unit testing
- [ ] Section 13 - Realtime engine with Socket.IO
- [ ] Section 14 - Deploy service with Docker
- [ ] Section 15 - Distributed tracing
- [ ] Section 16 - Final section - gRPC
