# Sử dụng Golang image chính thức
FROM golang:1.20-alpine

# Đặt thư mục làm việc
WORKDIR /app

# Sao chép file go.mod và go.sum (quản lý dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Sao chép toàn bộ mã nguồn vào container
COPY . .

# Biên dịch ứng dụng
RUN go build -o main .

# Mở port 8080 cho ứng dụng
EXPOSE 8080

# Lệnh chạy ứng dụng
CMD ["./main"]