# Hub Service API

Một backend service được xây dựng bằng Go và Gin framework, cung cấp các RESTful API cho các hoạt động của một hub trung tâm, bao gồm quản lý người dùng, xác thực và các tính năng khác.

## ✨ Tính năng

-   **Xác thực người dùng**: Đăng ký, Đăng nhập sử dụng JWT (JSON Web Tokens).
-   **Quản lý người dùng**: Các API theo chuẩn CRUD (Create, Read, Update, Delete) cho module người dùng.
-   **Kiến trúc Layered**: Phân tách rõ ràng giữa các lớp Transport, Business, và Storage.
-   **Middleware**: Tích hợp sẵn middleware cho logging, phục hồi (recovery) và xác thực (authentication).
-   **Tài liệu API**: Tự động sinh tài liệu API với Swagger.
-   **Quản lý cấu hình**: Dễ dàng quản lý cấu hình môi trường qua file `.env`.

## 🏗️ Kiến trúc

Dự án tuân theo kiến trúc phân lớp (Layered Architecture), lấy cảm hứng từ Clean Architecture để đảm bảo code được rõ ràng, dễ bảo trì và mở rộng.

```mermaid
graph TD
    subgraph "Client"
        A[User/Client Application]
    end

    subgraph "Hub Service"
        B(Gin Router)
        C{Middleware}
        D[Transport Layer<br/>(Handlers)]
        E[Business Logic Layer<br/>(Biz)]
        F[Storage Layer<br/>(Storage)]
        G((MongoDB))
    end

    A --> B
    B --> C
    C --> D
    D --> E
    E --> F
    F --> G
```

-   `transport`: Chịu trách nhiệm xử lý các request HTTP và response (sử dụng Gin handlers).
-   `biz` (Business Logic): Chứa toàn bộ logic nghiệp vụ của ứng dụng.
-   `storage`: Chịu trách nhiệm truy cập và thao tác với cơ sở dữ liệu (MongoDB).
-   `component`: Các thành phần cốt lõi, có thể tái sử dụng như kết nối database, provider cho token, hasher.
-   `model`: Định nghĩa các cấu trúc dữ liệu cho các lớp khác nhau.

## ⚙️ Luồng hoạt động chính

### 1. Luồng đăng nhập và tạo Access Token

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant UserHandler as Handler
    participant UserBiz as Biz
    participant TokenProvider as Provider
    participant Database

    Client->>+Router: POST /api/users/login (email, password)
    Router->>+UserHandler: Login()
    UserHandler->>+UserBiz: biz.Login(email, password)
    UserBiz->>+Database: GetUserByEmail(email)
    Database-->>-UserBiz: User record
    UserBiz->>UserBiz: Verify password
    UserBiz->>+TokenProvider: provider.Generate(payload)
    TokenProvider-->>-UserBiz: Access Token
    UserBiz-->>-UserHandler: LoginResponse (with token)
    UserHandler-->>-Router: JSON Response
    Router-->>-Client: 200 OK (gồm access_token)
```

### 2. Luồng xác thực Access Token khi gọi API

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant AuthMiddleware as Middleware
    participant ProtectedHandler as Handler

    Client->>+Router: GET /api/users (Header: Auth)
    Router->>+Middleware: Run AuthMiddleware
    Middleware->>Middleware: Validate JWT Token

    alt Token hợp lệ
        Middleware-->>+Handler: c.Next()
        Handler->>Handler: Process request
        Handler-->>-Router: JSON Response
        Router-->>-Client: 200 OK (data)
    else Token không hợp lệ
        Middleware-->>-Router: Abort with 401
        Router-->>-Client: 401 Unauthorized
    end
```

## 🚀 Bắt đầu

### Yêu cầu

-   [Go](https://golang.org/dl/) (phiên bản 1.23.1 hoặc mới hơn)
-   [MongoDB](https://www.mongodb.com/try/download/community)
-   [swag-cli](https://github.com/swaggo/swag)

### Cài đặt & Cấu hình

1.  **Clone repository:**

    ```bash
    git clone https://your-repository-url.com/hub-service.git
    cd hub-service
    ```

2.  **Tải dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Cấu hình môi trường:**
    Tạo một file `.env` từ file mẫu `env.txt` và điền các thông tin cần thiết.

    ```bash
    cp env.txt .env
    ```

    Nội dung file `.env`:

    ```env
    # MongoDB Configuration
    SYSTEM_MONGODB_URI="mongodb+srv://user:password@cluster.mongodb.net/your_db_name"
    SYSTEM_APP_NAME="your_db_name"

    # System Configuration
    SYSTEM_SECRET_KEY="a-very-strong-and-secret-key"

    # Server Configuration
    PORT="8080"
    ```

    -   `SYSTEM_MONGODB_URI`: Chuỗi kết nối đến MongoDB Atlas hoặc local.
    -   `SYSTEM_APP_NAME`: Tên database của bạn.
    -   `SYSTEM_SECRET_KEY`: Một chuỗi bí mật mạnh để tạo và xác thực JWT.
    -   `PORT`: Cổng mà server sẽ chạy.

### Chạy ứng dụng

-   **Chạy server:**

    ```bash
    go run main.go
    ```

    Server sẽ khởi động và chạy tại `http://localhost:8080` (hoặc cổng bạn đã cấu hình).

-   **Kiểm tra Health Check:**
    ```bash
    curl http://localhost:8080/health
    ```
    Response mong muốn:
    ```json
    {
    	"status": "healthy",
    	"message": "All services are running",
    	"mongodb": "connected"
    }
    ```

## 📚 Tài liệu API (Swagger)

Tài liệu API được tự động sinh bằng Swagger.

1.  **Generate tài liệu:**
    Sau khi thêm hoặc sửa đổi các API và các comment chú thích, chạy lệnh sau ở thư mục gốc của dự án:

    ```bash
    swag init
    ```

    Lệnh này sẽ cập nhật các file trong thư mục `/docs`.

2.  **Xem tài liệu:**
    Khởi động server và truy cập vào URL sau trên trình duyệt:
    `http://localhost:8080/swagger/index.html`

### Sử dụng API cần xác thực

-   Đầu tiên, gọi API `POST /api/users/login` để nhận `access_token`.
-   Đối với các API được bảo vệ, bạn cần thêm token vào header `Authorization`.
    ```
    Authorization: Bearer <your_access_token>
    ```
-   Bạn có thể sử dụng nút **"Authorize"** trên giao diện Swagger để lưu token và tự động đính kèm vào các request.

## 📁 Cấu trúc thư mục

```
.
├── common/           # Các hàm, hằng số, và lỗi chung
├── component/        # Các thành phần cốt lõi (database, auth, hasher)
├── docs/             # File tài liệu được sinh bởi Swagger
├── middleware/       # Gin middleware (CORS, recover, services router)
├── module/           # Các module nghiệp vụ chính
│   └── user/
│       ├── biz/      # Business logic cho user
│       ├── model/    # Data models cho user
│       ├── storage/  # Tương tác database cho user
│       └── transport/# API handlers và routes cho user
├── .env.example      # File mẫu cho biến môi trường
├── go.mod            # Quản lý dependencies
├── main.go           # Entry point của ứng dụng
└── README.md         # File này
```

Chúc bạn coding vui vẻ!
