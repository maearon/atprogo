# AtProGo with Neon PostgreSQL

A Go implementation of a decentralized social network protocol inspired by Bluesky's AT Protocol, using Neon PostgreSQL for data storage.

## Overview

AtProGo is a lightweight, modular implementation of a decentralized social network protocol using Go's standard library. It follows the architectural patterns of Bluesky's AT Protocol, focusing on:

- Decentralized identity and data ownership
- Content-addressed data structures
- Federation between independent servers
- Modular services architecture

## Architecture

The project is structured as a set of microservices:

- **api-gateway**: Routes requests to appropriate services
- **auth-service**: Handles authentication and identity
- **pds**: Personal Data Server for user data storage
- **bgs**: Big Graph Service for social graph operations

## Database Schema

The application uses the following tables:

1. **users**: Stores user information for authentication
2. **repositories**: Stores repository metadata
3. **commits**: Stores repository commits
4. **documents**: Stores repository documents
5. **follows**: Stores follow relationships

## Getting Started

\`\`\`bash
# Clone the repository
git clone https://github.com/yourusername/atprogo.git

# Set environment variables
export DATABASE_URL="your-neon-postgres-connection-string"

# Run the services
cd atprogo
go run cmd/api-gateway/main.go
\`\`\`

## API Endpoints

### Auth Service (port 8081)

- `POST /register`: Register a new user
- `POST /login`: Login a user

### PDS Service (port 8082)

- `POST /posts/create`: Create a new post
- `GET /posts/get?did={did}`: Get posts for a user

### BGS Service (port 8083)

- `POST /follow`: Follow a user
- `POST /unfollow`: Unfollow a user
- `GET /followers?did={did}`: Get followers for a user
- `GET /following?did={did}`: Get users a user is following

### API Gateway (port 8080)

- `*` /auth/*: Routes to Auth Service
- `*` /pds/*: Routes to PDS Service
- `*` /bgs/*: Routes to BGS Service

## License

MIT
