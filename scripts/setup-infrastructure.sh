#!/bin/bash
# Server Setup Script for Hub Service Infrastructure
# Run this ONCE on the server to setup Redis and Kafka

set -e

echo "=== Hub Service Infrastructure Setup ==="

# 1. Create shared network
echo "Creating Docker network..."
docker network create app-network 2>/dev/null || echo "Network 'app-network' already exists"

# 2. Start Redis
echo "Starting Redis..."
docker rm -f redis 2>/dev/null || true
docker run -d \
  --name redis \
  --network app-network \
  --restart always \
  -p 6379:6379 \
  -v redis_data:/data \
  redis:alpine

# 3. Start Kafka
echo "Starting Kafka..."
docker rm -f kafka 2>/dev/null || true
docker run -d \
  --name kafka \
  --network app-network \
  --restart always \
  -p 9092:9092 \
  -e KAFKA_NODE_ID=1 \
  -e KAFKA_PROCESS_ROLES=broker,controller \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092 \
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@kafka:9093 \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1 \
  -e CLUSTER_ID=MkU3OEVBNTcwNTJENDM2Qk \
  apache/kafka:3.7.0

echo ""
echo "=== Setup Complete ==="
echo "Redis: running on port 6379"
echo "Kafka: running on port 9092"
echo ""
echo "Update your hub-service.env with:"
echo "  REDIS_HOST=redis"
echo "  KAFKA_BROKERS=kafka:9092"
echo ""
echo "Now push code to main branch to deploy hub-service."
