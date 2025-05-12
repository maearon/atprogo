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

```
Chào mừng bạn đến với thế giới Go! 👋 Với background từ Ruby on Rails, bạn chắc chắn sẽ thấy kiến trúc của Go có phần "kỳ lạ" hoặc "bất tiện" lúc đầu — vì Go không có framework "full-stack" như Rails mà thiên về kiến trúc clean architecture hoặc modular.

Dưới đây là bản giải thích sơ đồ cây thư mục mà bạn đang thấy, theo kiểu dễ hiểu cho người từ Rails chuyển sang:

🏗️ Tổng quan thư mục (phía Go)
bash
Sao chép
Chỉnh sửa
.
├── cmd/
│   └── api-gateway/          # Entry point của app (kiểu như app/controllers trong Rails)
│       └── main.go           # Tương đương file boot hoặc application.rb
│
├── pkg/                      # Chứa các logic tách biệt, kiểu như lib/ hoặc concerns/
│   ├── auth/                 # Logic xác thực, giống concerns/authentication
│   ├── bgs/                  # Big Graph Service – quản lý follow
│   ├── db/                   # Kết nối DB hoặc query helper
│   ├── identity/             # Quản lý DID (Decentralized ID)
│   ├── lexicon/              # Chứa schema hoặc mô tả API
│   ├── pds/                  # Personal Data Service – CRUD tài liệu
│   ├── repo/                 # Lớp làm việc với dữ liệu, kiểu như ActiveRecord
│   └── xrpc/                 # XRPC protocol, giống như API layer
│
├── services/                 # Các service tương ứng với từng domain (auth, pds, bgs)
│   ├── auth/
│   ├── bgs/
│   ├── pds/
│   └── api-gateway/          # API Gateway: định tuyến đến các service kia
│
├── docker-compose.yml        # Dùng để chạy tất cả service trong container
├── Dockerfile.txt            # Dockerfile (có thể dùng để build image cho service)
├── go.mod                    # Tương đương Gemfile
├── init-db.sql               # File init dữ liệu ban đầu
└── README.md
🔄 So sánh với Rails
Go App (kiểu AI tạo)	Tương đương trong Rails	Ghi chú
cmd/api-gateway	config/routes.rb + controller khởi chạy	Điểm bắt đầu của app
pkg/auth	app/models/concerns/auth	Xác thực
pkg/repo	app/models	Logic thao tác DB
pkg/pds, pkg/bgs	app/services/	Logic nghiệp vụ
services/pds	app/controllers/pds_controller.rb	Từng service riêng
go.mod	Gemfile	Khai báo dependencies
Dockerfile	Dockerfile của Rails	
init-db.sql	db/schema.rb hoặc seeds.rb	

📌 Mẹo cho người mới từ Rails:
Không có ActiveRecord – Bạn phải viết SQL bằng tay hoặc dùng ORM nhẹ như sqlx, gorm, hoặc pgx.

Không có callbacks/DSL – Go ưu tiên rõ ràng, đơn giản, và khai báo tường minh.

Chia nhỏ logic rất mạnh – Bạn có thể thấy hơi verbose, nhưng về sau dễ scale, test, maintain.

```
