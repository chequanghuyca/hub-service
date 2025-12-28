package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	brokers []string
	writers map[string]*kafka.Writer
	readers map[string]*kafka.Reader
}

type Message struct {
	Key   []byte
	Value []byte
}

func NewKafkaClient() (*KafkaClient, error) {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092"
	}

	brokers := strings.Split(brokersEnv, ",")

	// Validate connection by creating a temporary connection
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	conn.Close()

	log.Println("Kafka connection established successfully!")

	return &KafkaClient{
		brokers: brokers,
		writers: make(map[string]*kafka.Writer),
		readers: make(map[string]*kafka.Reader),
	}, nil
}

// GetWriter returns a writer for the specified topic, creating one if it doesn't exist
func (k *KafkaClient) GetWriter(topic string) *kafka.Writer {
	if writer, exists := k.writers[topic]; exists {
		return writer
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
	}

	k.writers[topic] = writer
	return writer
}

// Produce sends a message to the specified topic
func (k *KafkaClient) Produce(ctx context.Context, topic string, key, value []byte) error {
	writer := k.GetWriter(topic)

	msg := kafka.Message{
		Key:   key,
		Value: value,
	}

	return writer.WriteMessages(ctx, msg)
}

// ProduceMultiple sends multiple messages to the specified topic
func (k *KafkaClient) ProduceMultiple(ctx context.Context, topic string, messages []Message) error {
	writer := k.GetWriter(topic)

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, msg := range messages {
		kafkaMessages[i] = kafka.Message{
			Key:   msg.Key,
			Value: msg.Value,
		}
	}

	return writer.WriteMessages(ctx, kafkaMessages...)
}

// GetReader returns a reader for the specified topic and group, creating one if it doesn't exist
func (k *KafkaClient) GetReader(topic, groupID string) *kafka.Reader {
	key := fmt.Sprintf("%s:%s", topic, groupID)
	if reader, exists := k.readers[key]; exists {
		return reader
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        500 * time.Millisecond,
		CommitInterval: time.Second,
	})

	k.readers[key] = reader
	return reader
}

// Consume reads a message from the specified topic
func (k *KafkaClient) Consume(ctx context.Context, topic, groupID string) (*kafka.Message, error) {
	reader := k.GetReader(topic, groupID)
	msg, err := reader.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// StartConsumer starts a consumer goroutine that processes messages with the provided handler
func (k *KafkaClient) StartConsumer(ctx context.Context, topic, groupID string, handler func(msg *kafka.Message) error) {
	reader := k.GetReader(topic, groupID)

	go func() {
		log.Printf("Starting Kafka consumer for topic: %s, group: %s", topic, groupID)
		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopping Kafka consumer for topic: %s", topic)
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return // Context cancelled
					}
					log.Printf("Error reading message from Kafka: %v", err)
					continue
				}

				if err := handler(&msg); err != nil {
					log.Printf("Error processing message: %v", err)
					// Could implement retry logic here
				}
			}
		}
	}()
}

// CreateTopic creates a topic with the specified configuration
func (k *KafkaClient) CreateTopic(topic string, partitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}

// HealthCheck verifies Kafka connection is alive
func (k *KafkaClient) HealthCheck() error {
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

// Close closes all writers and readers
func (k *KafkaClient) Close() error {
	var lastErr error

	for _, writer := range k.writers {
		if err := writer.Close(); err != nil {
			lastErr = err
		}
	}

	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// GetBrokers returns the list of Kafka brokers
func (k *KafkaClient) GetBrokers() []string {
	return k.brokers
}
