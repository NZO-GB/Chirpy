Chirpy

A Twitter-inspired REST API written in Go.

Chirpy is a backend service that implements the core functionality of a social media platform: user registration, authentication, posting short messages ("chirps"), and account management. The project was built to explore backend development in Go while following common REST API design patterns and modern authentication practices.

Features
User registration and authentication
Secure password hashing with Argon2id
JWT-based authentication
Refresh token support
Token revocation (logout)
Create, retrieve, and delete chirps
Update user credentials
Filter chirps by author
Sort chirps in ascending or descending order
PostgreSQL persistence
SQL migrations with Goose
Type-safe database queries using sqlc
Webhook endpoint for Chirpy Red upgrades
Tech Stack
Technology	Purpose
Go	Backend API
PostgreSQL	Relational database
sqlc	Type-safe SQL query generation
Goose	Database migrations
Argon2id	Password hashing
JWT	Authentication
godotenv	Environment variable management
Project Structure
.
├── internal/
│   └── database/      # sqlc-generated database layer
├── sql/
│   ├── schema/        # Goose migrations
│   └── queries/       # SQL used by sqlc
├── main.go
├── go.mod
└── README.md
Configuration

Create a .env file in the project root.

DB_URL=postgres://username:password@localhost:5432/chirpy?sslmode=disable
SECRET=your_jwt_secret
POLKA_KEY=your_polka_webhook_key
PLATFORM=dev
Environment Variables
Variable	Description
DB_URL	PostgreSQL connection string
SECRET	Secret used to sign JWTs
POLKA_KEY	API key used to validate Polka webhook requests
PLATFORM	Application environment (for example dev)
Running the Project
1. Clone the repository
git clone <repository-url>
cd chirpy
2. Create a PostgreSQL database

Create an empty PostgreSQL database and configure DB_URL accordingly.

3. Run database migrations
goose postgres "$DB_URL" up
4. Start the server
go run .

The server listens on:

http://localhost:8080
Authentication Flow
Register a user.
Log in with your credentials.
Receive an access token and a refresh token.
Include the access token in authenticated requests:
Authorization: Bearer <access_token>
When the access token expires, request a new one using the refresh token.
Revoke the refresh token when logging out.
API Endpoints
Authentication
Method	Endpoint	Description
POST	/api/users	Register a new user
POST	/api/login	Authenticate a user
POST	/api/refresh	Generate a new access token
POST	/api/revoke	Revoke a refresh token
PUT	/api/users	Update user credentials
Chirps
Method	Endpoint	Description
POST	/api/chirps	Create a chirp
GET	/api/chirps	Retrieve all chirps
GET	/api/chirps/{chirpID}	Retrieve a single chirp
DELETE	/api/chirps/{chirpID}	Delete a chirp
Webhooks
Method	Endpoint	Description
POST	/api/polka/webhooks	Receive webhook events from Polka and upgrade eligible users to Chirpy Red
Health Check
Method	Endpoint	Description
GET	/api/healthz	Verify that the server is running
Database

Chirpy uses PostgreSQL as its primary datastore.

Database schema changes are managed with Goose migrations, while sqlc generates type-safe Go code directly from SQL queries, allowing the project to retain the expressiveness of handwritten SQL without sacrificing compile-time safety.

Future Improvements

Possible future enhancements include:

Pagination for chirp feeds
User profiles
Likes and replies
Follower relationships
Rate limiting
Comprehensive automated tests
Docker Compose deployment
OpenAPI (Swagger) documentation
Acknowledgements

This project was originally developed as part of the Boot.dev backend curriculum and later expanded into a standalone REST API project.