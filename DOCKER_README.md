# Docker Deployment Guide

## 🐳 Quick Start

### Prerequisites

-   Docker Engine 20.10+
-   Docker Compose 2.0+
-   At least 2GB RAM available

### 1. Development Environment

```bash
# Clone the repository
git clone <your-repo-url>
cd hub-service

# Copy environment file
cp env.example .env

# Edit environment variables
nano .env

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f hub-service
```

### 2. Production Environment

```bash
# Copy environment file
cp env.example .env

# Edit environment variables for production
nano .env

# Start production services
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f hub-service
```

## 📋 Services

### Development (`docker-compose.yml`)

-   **Hub Service**: `http://localhost:8080`
-   **MongoDB**: `localhost:27017`
-   **Mongo Express**: `http://localhost:8081` (admin/password123)

### Production (`docker-compose.prod.yml`)

-   **Hub Service**: `http://localhost:8080` (localhost only)
-   **MongoDB**: `localhost:27017` (localhost only)
-   **Nginx**: `http://localhost:80` (optional)

## 🔧 Environment Variables

### Required Variables

```bash
SYSTEM_SECRET_KEY=your-super-secret-key-here
SYSTEM_MONGODB_URI=mongodb://admin:password123@mongodb:27017/hub_service?authSource=admin
```

### Optional Variables

```bash
# Email Configuration
SYSTEM_EMAIL=your-email@gmail.com
SYSTEM_EMAIL_HOST=smtp.gmail.com
SYSTEM_EMAIL_PORT=587

# DeepL Configuration
SYSTEM_DEEPL_API_KEY=your-deepl-api-key
SYSTEM_DEEPL_BASE_URL=https://api-free.deepl.com/v2

# Portfolio URLs
BASE_URL_PORTFOLIO=https://your-portfolio.com
BASE_URL_LOCAL=http://localhost:3000
```

## 🚀 Deployment Commands

### Development

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# Rebuild and start
docker-compose up -d --build

# View logs
docker-compose logs -f

# Access MongoDB shell
docker-compose exec mongodb mongosh -u admin -p password123
```

### Production

```bash
# Start production services
docker-compose -f docker-compose.prod.yml up -d

# Stop production services
docker-compose -f docker-compose.prod.yml down

# Rebuild and start production
docker-compose -f docker-compose.prod.yml up -d --build

# View production logs
docker-compose -f docker-compose.prod.yml logs -f
```

## 🔍 Health Checks

### Application Health

```bash
curl http://localhost:8080/health
```

### MongoDB Health

```bash
docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

## 📊 Monitoring

### View Container Status

```bash
docker-compose ps
```

### Resource Usage

```bash
docker stats
```

### Logs

```bash
# All services
docker-compose logs

# Specific service
docker-compose logs hub-service

# Follow logs
docker-compose logs -f hub-service
```

## 🛠️ Troubleshooting

### Common Issues

1. **Port already in use**

    ```bash
    # Check what's using the port
    lsof -i :8080

    # Kill the process or change port in docker-compose.yml
    ```

2. **MongoDB connection failed**

    ```bash
    # Check MongoDB logs
    docker-compose logs mongodb

    # Restart MongoDB
    docker-compose restart mongodb
    ```

3. **Application won't start**

    ```bash
    # Check application logs
    docker-compose logs hub-service

    # Check environment variables
    docker-compose exec hub-service env | grep SYSTEM
    ```

### Reset Everything

```bash
# Stop and remove all containers, networks, volumes
docker-compose down -v

# Remove all images
docker-compose down --rmi all

# Start fresh
docker-compose up -d --build
```

## 🔐 Security Notes

### Production Security

1. **Change default passwords** in `.env` file
2. **Use strong SYSTEM_SECRET_KEY**
3. **Enable SSL/TLS** with Nginx
4. **Restrict network access** (already configured in prod)
5. **Regular security updates**

### Default Credentials

-   **MongoDB**: admin/password123
-   **Mongo Express**: admin/password123
-   **Default Super Admin**: admin@hubservice.com/password

## 📝 API Documentation

Once deployed, access Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## 🗄️ Database Management

### Access MongoDB

```bash
# Via Docker
docker-compose exec mongodb mongosh -u admin -p password123

# Via Mongo Express (Development)
http://localhost:8081
```

### Backup Database

```bash
# Create backup
docker-compose exec mongodb mongodump --out /data/backup

# Copy backup to host
docker cp hub-service-mongodb:/data/backup ./backup
```

### Restore Database

```bash
# Copy backup to container
docker cp ./backup hub-service-mongodb:/data/

# Restore
docker-compose exec mongodb mongorestore /data/backup
```

## 🔄 Updates

### Update Application

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build
```

### Update Dependencies

```bash
# Rebuild with no cache
docker-compose build --no-cache

# Restart services
docker-compose up -d
```
