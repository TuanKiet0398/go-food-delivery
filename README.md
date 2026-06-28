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

## Roadmap

- [x] Section 02 - UI to Database
- [ ] Section 03 - GORM
- [ ] Section 04 - Simple clean architecture
- [ ] Section 05 - Error handling and UID
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
