# FitQuest Backend

GraphQL API backend for FitQuest — a fitness exercise reference and workout tracking app.

Built with Go, [gqlgen](https://gqlgen.com/), MongoDB, and JWT authentication.

## Tech Stack

- **Go 1.26** — runtime
- **gqlgen** — GraphQL server code generation
- **MongoDB Atlas** — database
- **JWT** — authentication (HS256 tokens)
- **bcrypt** — password hashing

## Project Structure

```
.
├── auth/               # HTTP middleware + context helpers for JWT auth
├── database/           # MongoDB connection and models
├── graph/              # GraphQL schema, resolvers, generated code
│   ├── model/          # Generated GraphQL types
│   ├── users/          # User mutations + schema
│   │   ├── users.graphqls
│   │   └── users.resolvers.go
│   ├── generated.go    # Generated runtime (do not edit)
│   ├── schema.graphqls # Exercise GraphQL schema
│   ├── schema.resolvers.go  # Exercise query resolvers
│   ├── helpers.go      # Shared utility functions
│   ├── mapping.go      # Mongo-to-GQL model mappers
│   └── resolver.go     # Root Resolver struct
├── jwt/                # JWT token generation and parsing
├── server.go           # Entry point
├── gqlgen.yml          # gqlgen configuration
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

Pass the token as a `Bearer` header:

```json
{ "Authorization": "Bearer <token>" }
```

## Code Generation

If you modify `graph/schema.graphqls`, regenerate the code:

```bash
go run github.com/99designs/gqlgen generate
```
