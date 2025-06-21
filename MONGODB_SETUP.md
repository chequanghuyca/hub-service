# MongoDB Setup Guide

## Cấu hình MongoDB cho Hub Service

### 1. Environment Variables

Tạo file `.env` trong thư mục gốc với các biến sau:

```env
# MongoDB Configuration
SYSTEM_MONGODB_URI=mongodb+srv://hub-service:<db_password>@hocmongo.ywb2dbq.mongodb.net/?retryWrites=true&w=majority&appName=hocmongo
SYSTEM_APP_NAME=hub_service_db

# System Configuration
SYSTEM_SECRET_KEY=your-secret-key-here

# Server Configuration
PORT=8080
```

### 2. Thay thế password

Thay thế `<db_password>` trong `SYSTEM_MONGODB_URI` bằng password thực tế của bạn.

### 3. Chạy ứng dụng

```bash
go run main.go
```

### 4. Test connection

Truy cập endpoint health check:

```
GET http://localhost:8080/health
```

Response thành công:

```json
{
	"status": "healthy",
	"message": "All services are running",
	"mongodb": "connected"
}
```

### 5. Cấu trúc Database

-   **Database**: `hub_service_db`
-   **Collections**:
    -   `users` - Lưu trữ thông tin người dùng
    -   Các collections khác sẽ được tạo tự động khi cần

### 6. Kiến trúc Database Component

```
main.go → database.NewDatabase() → MongoDB Connection
    ↓
AppContext.GetDatabase() → Database Component
    ↓
User Storage → MongoDB Operations
```

### 7. Best Practices

1. **Separation of Concerns**: Database logic được tách riêng trong `component/database/`
2. **Security**: Không commit file `.env` vào git
3. **Connection Pooling**: MongoDB driver tự động quản lý connection pool
4. **Error Handling**: Luôn kiểm tra lỗi khi thực hiện database operations
5. **Indexing**: Tạo index cho các trường thường query
6. **Validation**: Sử dụng BSON tags để validate dữ liệu

### 8. Troubleshooting

-   **Connection failed**: Kiểm tra MongoDB URI và network connection
-   **Authentication failed**: Kiểm tra username/password
-   **Database not found**: Database sẽ được tạo tự động khi có dữ liệu đầu tiên
