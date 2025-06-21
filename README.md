# Hub Service API

## ✨ Tính năng

-   **Xác thực người dùng**: Đăng ký, Đăng nhập sử dụng JWT (JSON Web Tokens).
-   **Quản lý người dùng**: Các API theo chuẩn CRUD (Create, Read, Update, Delete) cho module người dùng.
-   **Kiến trúc Layered**: Phân tách rõ ràng giữa các lớp Transport, Business, và Storage.
-   **Middleware**: Tích hợp sẵn middleware cho logging, phục hồi (recovery) và xác thực (authentication).
-   **Tài liệu API**: Tự động sinh tài liệu API với Swagger.
-   **Quản lý cấu hình**: Dễ dàng quản lý cấu hình môi trường qua file `.env`.

## 🏗️ Kiến trúc Tổng quan

```mermaid
graph TD;
    subgraph "Client"
        A["User/Client Application"];
    end

    subgraph "Hub Service"
        B("Gin Router");
        C{"Middleware"};
        D["Transport Layer<br>(Handlers)"];
        E["Business Logic Layer<br>(Biz)"];
        F["Storage Layer<br>(Storage)"];
        G(("MongoDB"));
    end

    A --> B;
    B --> C;
    C --> D;
    D --> E;
    E --> F;
    F --> G;
```

## ⚙️ Luồng hoạt động chính

### 1. Luồng đăng nhập và tạo Access Token

Đây là quá trình người dùng cung cấp thông tin xác thực (email và password) để nhận về một `access_token`. Token này giống như một chiếc chìa khóa tạm thời để truy cập các tài nguyên khác.

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

Khi đã có `access_token`, người dùng sẽ đính kèm nó vào header `Authorization` của mỗi request đến các API cần xác thực. Middleware sẽ kiểm tra tính hợp lệ của token trước khi cho phép request đi tiếp.

```mermaid
sequenceDiagram
    participant Client
    participant Router
    participant AuthMiddleware
    participant ProtectedHandler

    Client->>+Router: GET /api/users (Header: Authorization: Bearer <token>)
    Router->>+AuthMiddleware: Run AuthMiddleware
    AuthMiddleware->>AuthMiddleware: Validate JWT Token

    alt Token hợp lệ
        AuthMiddleware->>+ProtectedHandler: c.Next()
        ProtectedHandler-->>-AuthMiddleware: Process request & return
    else Token không hợp lệ
        AuthMiddleware-->>Router: Abort with 401
    end

    AuthMiddleware-->>-Router: Pass control back
    Router-->>-Client: Final Response
```

### 3. Luồng Làm Mới Token (Refresh Token)

Để cải thiện trải nghiệm người dùng, thay vì bắt họ đăng nhập lại mỗi khi `access_token` (ngắn hạn) hết hạn, hệ thống sử dụng một `refresh_token` (dài hạn) để lấy một cặp token mới.

-   **Access Token**: Thời gian sống ngắn (ví dụ: 15 phút), dùng để truy cập tài nguyên.
-   **Refresh Token**: Thời gian sống dài (ví dụ: 30 ngày), chỉ dùng để lấy `access_token` mới.

Dưới đây là sơ đồ mô tả luồng hoạt động khi `access_token` hết hạn:

```mermaid
sequenceDiagram
    participant Client
    participant Server

    %% 1. Access Token Expires
    Note over Client,Server: Access Token đã hết hạn

    Client->>Server: GET /api/users (Header: Bearer access_token hết hạn)
    Server->>Server: Xác thực Access Token -> Thất bại (Hết hạn)
    Server-->>Client: 401 Unauthorized Error

    %% 2. Refresh Token Flow
    Client->>Server: POST /api/users/refresh (Body: { refresh_token })
    Server->>Server: Xác thực Refresh Token
    alt Refresh Token Hợp lệ
        Server->>Server: Tạo Access Token MỚI & Refresh Token MỚI
        Server-->>Client: { new_access_token, new_refresh_token }
        Client->>Client: Cập nhật, lưu trữ 2 token MỚI
    else Refresh Token Không Hợp lệ
        Server-->>Client: 401 Unauthorized Error (Yêu cầu đăng nhập lại)
    end

    %% 3. Retry Original Request
    Client->>Server: GET /api/users (Header: Bearer new_access_token)
    Server->>Server: Xác thực Access Token mới -> Thành công
    Server-->>Client: Dữ liệu người dùng
```
