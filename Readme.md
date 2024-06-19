# Go Fiber Backend Template

Welcome to the Go Fiber Backend Template! This template helps you kickstart backend development with Go and Fiber, following Clean Architecture principles with Repository pattern and Test-Driven Development (TDD).

## Features

- **Clean Architecture**: Separation of concerns with distinct layers.
- **Fiber Framework**: Fast and efficient web framework for Go.
- **JWT Authentication**: Secure token-based authentication.
- **Dependency Injection**: Managed with Google Wire.
- **PostgreSQL**: Integrated database support.
- **Validation**: Request validation using `go-playground/validator`.
- **Caching**: Integrated with Redis.
- **Docker**: Containerization with Docker and Docker Compose.
- **Argon2**: Strong algorithm for password hashing.

## Project Structure

```
├── applications
│   ├── cache
│   ├── security
│   ├── validation
├── commons
├── domains
│   ├── entity
│   ├── repository
├── infrastructures
│   ├── database
│   ├── redis
│   ├── router
├── interfaces
│   ├── http
│   ├── middleware
├── migrations
├── tests
├── .env
├── .gitignore
├── Makefile
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

## Getting Started
### Prerequisites
- Go 1.18 or later
- Docker
- Docker Compose

## Installation
### 1. Clone the repository:

```bash
git clone https://github.com/wisle25/be-template.git
cd be-template
```

### 2. Install dependencies:

```bash
go mod tidy
```

### 3. Setup Environment Variables

```dotenv
# DATABASE and CACHING
POSTGRES_HOST=your_postgres_host_here
POSTGRES_PORT=your_postgres_port_here
POSTGRES_USER=your_postgres_user_here
POSTGRES_PASSWORD=your_postgres_password_here
POSTGRES_DB=your_postgres_db_here
POSTGRES_DB_TEST=your_postgres_db_test_here

PGADMIN_DEFAULT_EMAIL=your_pgadmin_email_here
PGADMIN_DEFAULT_PASSWORD=your_pgadmin_password_here

REDIS_URL=your_redis_url_here

# SERVER
APP_ENV=dev # Change this to "prod" for production
PORT=8000
CLIENT_ORIGIN=http://localhost:3000

# JWT
ACCESS_TOKEN_PRIVATE_KEY=your_access_token_private_key_here
ACCESS_TOKEN_PUBLIC_KEY=your_access_token_public_key_here
ACCESS_TOKEN_EXPIRED_IN=5m
ACCESS_TOKEN_MAXAGE=5

REFRESH_TOKEN_PRIVATE_KEY=your_refresh_token_private_key_here
REFRESH_TOKEN_PUBLIC_KEY=your_refresh_token_public_key_here
REFRESH_TOKEN_EXPIRED_IN=60m
REFRESH_TOKEN_MAXAGE=60
```

### 4. Compose docker
```bash
docker-compose up -d
```

### 5. Test
```bash
go test ./...
```

### 6. Build and run
```bash
go build -ldflags "-s -w"
./your-app.exe
```

## API Endpoints
### 1. User Registration
- Endpoint: POST /users
- Payload:
```json
{
    "username": "string",
    "email": "string",
    "password": "string", 
    "confirmPassword": "string"
}
```

- Response:

```json
{
    "status": "success",
    "data": "new registered user id"
}
```
### 2. User Login
- Endpoint: POST /auths
- Payload:

```json
{
    "identity": "Identity can be username or email",
    "password": "string"
}
```
- Response:

``` json
{
    "status": "success",
    "message": "Successfully logged in!"
}
```

### 3. Refresh Token
- Endpoint: PUT /auths
- Payload:
```
No payload needed, the refresh token is retrieved from cookies.
```

- Response:

``` json
{
    "status": "success"
}
```

### 4. Logout
- Endpoint: DELETE /auths
- Payload:
```
No payload needed, the tokens are retrieved from cookies.
```
- Response:
```json
{
    "status": "success",
    "message": "Successfully logged out!"
}
```

## Contributing
Contributions are welcome! Please fork this repository and submit pull requests.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Contact
For any questions or suggestions, feel free to open an issue or contact me at my [email](handidwic1225@gmail.com).