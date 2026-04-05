# go-app-shared

Shared Go module containing the DTOs used for async inter-microservice communication across the go-app platform.

## Purpose

This module acts as a contract layer between microservices. Each DTO defines the shape of a message exchanged over a messaging system (RabbitMQ, and potentially others in the future). Keeping them in a single shared module ensures all producers and consumers stay in sync.

## Package Structure

```
shared/
└── messaging/
    └── rabbitmq/
        └── dtos/
            ├── welcome_email.go
            └── user_logged_in.go
```

New messaging systems (e.g. Kafka) should follow the same layout under `messaging/<system>/dtos/`.

## DTOs

### RabbitMQ — `messaging/rabbitmq/dtos`

| Struct | Exchange | Routing Key | Published by | Consumed by |
|---|---|---|---|---|
| `WelcomeEmail` | `auth.events` | `user.created` | auth | email |
| `UserLoggedIn` | `auth.events` | `user.logged_in` | auth | broadcasting |

#### `WelcomeEmail`

```go
type WelcomeEmail struct {
    Email           string `json:"email"`
    Name            string `json:"name"`
    VerificationURL string `json:"verification_url"`
}
```

#### `UserLoggedIn`

```go
type UserLoggedIn struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}
```

## Usage

Import the module using its Go module path:

```go
import dtos "github.com/guille1988/go-app-shared/messaging/rabbitmq/dtos"
```

## Extending

To add a new messaging system, create the corresponding directory and package:

```
messaging/
└── kafka/
    └── dtos/
        └── some_event.go   // package dtos
```

Follow the same struct + JSON tag conventions used in the RabbitMQ DTOs.
