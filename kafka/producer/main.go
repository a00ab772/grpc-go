package main

import (
	"fmt"
	"log"

	pb "example.com/m/kafka/proto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

func main() {
	// 1. Schema Registry Client
	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))
	if err != nil {
		log.Fatalf("Failed to create SR client: %s", err)
	}

	// 2. Configure Protobuf Serializer
	// Removing the manual TopicNameStrategy assignment to fix the build error.
	serConfig := protobuf.NewSerializerConfig()
	ser, err := protobuf.NewSerializer(srClient, serde.ValueSerde, serConfig)
	if err != nil {
		log.Fatalf("Failed to create serializer: %s", err)
	}

	// 3. Create Kafka Producer (Use 9093 as configured in docker-compose)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	topic := "user-updates-value"

	// 4. Produce Messages
	for i := 1; i <= 10; i++ {
		user := &pb.UserRequest{Id: int64(i)}

		payload, err := ser.Serialize(topic, user)
		if err != nil {
			log.Printf("Serialization error: %s", err)
			continue
		}

		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          payload,
		}, nil)

		if err != nil {
			log.Printf("Produce error: %s", err)
		}
	}

	p.Flush(15 * 1000)
	fmt.Println("Produced 10 messages successfully.")
}
