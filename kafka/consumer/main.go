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
	// Must match the topic used in the producer
	topic := "user-updates-value"
	// The subject used to fetch the schema from the registry
	subject := "user-updates-value"

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9093",
		"group.id":          "my-local-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{topic}, nil)

	// Initialize SR Client
	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))
	if err != nil {
		log.Fatalf("Failed to create SR client: %s", err)
	}

	// Initialize Deserializer
	deser, err := protobuf.NewDeserializer(srClient, serde.ValueSerde, protobuf.NewDeserializerConfig())
	if err != nil {
		log.Fatalf("Failed to create deserializer: %s", err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			user := &pb.UserRequest{}
			// DeserializeInto uses the registered schema for this subject
			err := deser.DeserializeInto(subject, msg.Value, user)
			if err != nil {
				log.Printf("Deserialization error: %s", err)
				continue
			}
			fmt.Printf("Received user ID: %d\n", user.Id)
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
