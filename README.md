# Go Project Starter Template

A production-ready Go project scaffold that provides a solid foundation for building scalable web applications using modern Go practices and popular libraries.

## Features

- 🚀 Modern project structure following Go best practices
- 🔒 JWT-based authentication
- 📝 Structured logging with rotation (Zap + Lumberjack)
- 🗄️ Database integration with GORM
- 🔄 Redis caching support
- ⚡ Dependency injection using Wire
- 🔧 YAML-based configuration
- 🧪 Testing setup with mocks
- 🛡️ Middleware support
- 🎯 Clean architecture pattern

## Project Structure

```
.
├── cmd/
│ └── api/ # Application entrypoints
├── config/ # Configuration files
├── internal/ # Private application code
│ ├── api/ # HTTP handlers
│ ├── middleware/ # HTTP middleware
│ ├── model/ # Domain models
│ ├── repository/ # Data access layer
│ ├── router/ # HTTP router setup
│ ├── service/ # Business logic
│ └── wire/ # Dependency injection
├── pkg/ # Public libraries
│ ├── cache/ # Caching utilities
│ ├── database/ # Database utilities
│ └── logger/ # Logging utilities
└── scripts/ # Build/deployment
```

## Technical Stack

- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Logging**: [Zap](https://github.com/uber-go/zap)
- **Caching**: [go-redis](https://github.com/redis/go-redis)
- **Authentication**: JWT using [golang-jwt](https://github.com/golang-jwt/jwt)
- **Dependency Injection**: [Wire](https://github.com/google/wire)

## Getting Started

### Prerequisites

- Go 1.23 or higher
- MySQL
- Redis

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/go-project-starter.git
cd go-project-starter
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example configuration:
```bash
cp config/config.example.yaml config/config.yaml
```

4. Update the configuration in `config/config.yaml` with your settings.

### Running the Application

1. Start the development server:
```bash
make run
# or
go run cmd/api/main.go
```

2. Run tests:
```bash
make test
# or
go test ./...
```

## Configuration

The application uses YAML-based configuration. Key configuration options include:

```yaml
server:
  port: 8080
  mode: development # or production

database:
  driver: mysql
  host: localhost
  port: 3306
  name: myapp
  user: root
  password: secret

redis:
  host: localhost
  port: 6379
  db: 0

jwt:
  secret: your-secret-key
  expiration: 24h

logging:
  level: info
  filename: app.log
  maxSize: 100
  maxBackups: 3
  maxAge: 28
```

## Project Layout Explanation

- `cmd/`: Contains the main applications of the project
- `internal/`: Contains the private application and library code
  - `api/`: HTTP handlers and route definitions
  - `middleware/`: HTTP middleware components
  - `model/`: Domain models and business logic interfaces
  - `repository/`: Data access layer implementations
  - `service/`: Business logic implementations
- `pkg/`: Libraries that can be used by external applications
- `config/`: Configuration files and templates

## Development

### Adding New Features

1. Define your domain models in `internal/model/`
2. Implement the repository interface in `internal/repository/`
3. Add business logic in `internal/service/`
4. Create HTTP handlers in `internal/api/`
5. Register routes in `internal/router/`

### Testing

The project includes examples for:
- Unit tests
- Integration tests
- Mock implementations

Run tests with coverage:
```bash
make test-coverage
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)