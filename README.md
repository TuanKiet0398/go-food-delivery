# food-delivery

A Go learning project using GORM with MySQL.

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
