# FitQuest Backend

GraphQL API backend for FitQuest — a fitness exercise reference and workout tracking app.

Built with Go, [gqlgen](https://gqlgen.com/), MongoDB, and JWT authentication. Resolvers depend on repository interfaces rather than the MongoDB driver directly, making them testable without a database connection.

## Tech Stack

- **Go 1.26** — runtime
- **gqlgen** — GraphQL server code generation
- **MongoDB Atlas** — database
- **JWT** — authentication (HS256 tokens)
- **bcrypt** — password hashing

## Project Structure

```
.
├── auth/                  # HTTP middleware + context helpers for JWT auth
├── database/              # MongoDB connection, models, and repository layer
│   ├── database.go        # DB connection
│   ├── models.go          # BSON models
│   └── repos.go           # Repository interfaces + Mongo implementations
├── graph/                 # GraphQL schema, resolvers, generated code
│   ├── generated.go       # Generated runtime (do not edit)
│   ├── graphqls/          # Schema files by domain
│   │   ├── exercises/
│   │   │   └── exercises.graphqls
│   │   └── users/
│   │       └── users.graphqls
│   ├── helper/            # Shared utility functions
│   │   └── helpers.go
│   ├── mapping/           # Mongo-to-GQL model mappers
│   │   ├── mapping.go
│   │   └── tests/
│   │       └── mapping_test.go
│   ├── model/             # Generated GraphQL types
│   │   └── models_gen.go
│   └── resolvers/         # Resolver implementations
│       ├── resolver.go    # Root Resolver struct (holds repo interfaces)
│       ├── exercises.resolvers.go
│       ├── users.resolvers.go
│       └── tests/
│           ├── exercises.resolvers_test.go
│           └── users.resolvers_test.go
├── jwt/                   # JWT token generation and parsing
├── server.go              # Entry point
├── gqlgen.yml             # gqlgen configuration
├── .golangci.yml          # Linter config
└── go.mod
```

## Prerequisites

- Go 1.26+
- A MongoDB Atlas cluster (or local MongoDB instance)

## Setup

### 1. Clone and install dependencies

```bash
git clone <repo-url>
cd fitquest-backend
go mod tidy
```

### 2. Configure environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Edit `.env`:

```
MONGODB_URI=mongodb+srv://<username>:<password>@<cluster>.mongodb.net/?appName=<app-name>
JWT_SECRET=<your-256-bit-random-secret>
PORT=5000
```

### 3. Seed the database

The `exercises` collection needs to exist in the `fitquest` database. There is no built-in seed script — import your exercise data manually (e.g. via MongoDB Compass, `mongosh`, or a custom script). Each document should follow this structure:

```json
{
  "_id": ObjectId,
  "name": "Bench Press",
  "category": "Strength",
  "mechanic": "Compound",
  "primary_muscles": ["Chest"],
  "secondary_muscles": ["Triceps", "Front Delts"]
}
```

### 4. Run the server

```bash
go run server.go
```

The GraphQL playground is available at `http://localhost:5000/`.

## GraphQL API

### Mutations

**Register** — Create a new user account:

```graphql
mutation {
  register(username: "alice", password: "securepass") {
    token
    user { id username }
  }
}
```

**Login** — Authenticate and receive a JWT:

```graphql
mutation {
  login(username: "alice", password: "securepass") {
    token
    user { id username }
  }
}
```

### Queries

**getExercises** — List all exercises (requires auth):

```graphql
query {
  getExercises {
    id
    name
    category
    mechanic
    primaryMuscles
    secondaryMuscles
  }
}
```

**getFilteredExercises** — List exercises filtered by muscle (requires auth):

```graphql
query {
  getFilteredExercises(where: { muscle: "Chest" }) {
    id
    name
    primaryMuscles
    secondaryMuscles
  }
}
```

Pass the token as a `Bearer` header:

```json
{ "Authorization": "Bearer <token>" }
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./graph/resolvers/tests/ -v
```

The resolver tests use repository interfaces to mock the database layer, so they run without a real MongoDB connection.

## Code Generation

If you modify any `.graphqls` schema files, regenerate the code:

```bash
go run github.com/99designs/gqlgen generate
```
