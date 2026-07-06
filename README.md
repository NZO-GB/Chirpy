# Chirpy

> A Twitter-inspired REST API built with Go, PostgreSQL, and JWT
> authentication.

Chirpy is a backend service that implements the core functionality of a
microblogging platform: user registration, authentication, posting short
messages ("chirps"), and account management.

## Features

-   User registration and login
-   JWT authentication with refresh tokens
-   Argon2id password hashing
-   Create, retrieve, filter, and delete chirps
-   Update user credentials
-   PostgreSQL persistence
-   Goose database migrations
-   sqlc-generated type-safe database access
-   Polka webhook integration

------------------------------------------------------------------------

## Tech Stack

  Technology   Purpose
  ------------ -------------------------------
  Go           Backend API
  PostgreSQL   Database
  sqlc         Type-safe SQL code generation
  Goose        Database migrations
  Argon2id     Password hashing
  JWT          Authentication
  godotenv     Environment configuration

------------------------------------------------------------------------

## Project Structure

``` text
.
├── internal/
│   └── database/
├── sql/
│   ├── queries/
│   └── schema/
├── main.go
├── go.mod
└── README.md
```

------------------------------------------------------------------------

## Configuration

Create a `.env` file:

``` env
DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
SECRET=your_jwt_secret
POLKA_KEY=your_polka_webhook_key
PLATFORM=dev
```

  Variable      Description
  ------------- --------------------------------------
  `DB_URL`      PostgreSQL connection string
  `SECRET`      JWT signing secret
  `POLKA_KEY`   Polka webhook API key
  `PLATFORM`    Application environment (e.g. `dev`)

------------------------------------------------------------------------

## Running the Project

### Prerequisites

-   Go
-   PostgreSQL
-   Goose

### Clone the repository

``` bash
git clone <repository-url>
cd chirpy
```

### Run the database migrations

``` bash
goose postgres "$DB_URL" up
```

### Start the server

``` bash
go run .
```

The API will be available at `http://localhost:8080`.

------------------------------------------------------------------------

## Authentication Flow

1.  Register a user.
2.  Log in.
3.  Receive an access token and refresh token.
4.  Use the access token:

``` http
Authorization: Bearer <access_token>
```

5.  Refresh the access token when it expires.
6.  Revoke the refresh token to log out.

------------------------------------------------------------------------

## API Endpoints

  Method   Endpoint                  Description
  -------- ------------------------- --------------------------
  POST     `/api/users`              Register a new user
  POST     `/api/login`              Authenticate a user
  POST     `/api/refresh`            Issue a new access token
  POST     `/api/revoke`             Revoke a refresh token
  PUT      `/api/users`              Update user credentials
  POST     `/api/chirps`             Create a chirp
  GET      `/api/chirps`             List chirps
  GET      `/api/chirps/{chirpID}`   Retrieve a chirp
  DELETE   `/api/chirps/{chirpID}`   Delete a chirp
  POST     `/api/polka/webhooks`     Handle Polka webhooks
  GET      `/api/healthz`            Health check

------------------------------------------------------------------------

## Future Improvements

-   Pagination
-   User profiles
-   Likes and replies
-   Follower relationships
-   Rate limiting
-   Automated tests
-   Docker Compose support
-   OpenAPI / Swagger documentation

------------------------------------------------------------------------

## Acknowledgements

Originally developed as part of the Boot.dev Backend Path and expanded
into a standalone portfolio project.
