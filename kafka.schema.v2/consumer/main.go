package main

import (
	"fmt"
	"log"

	pb "example.com/m/kafka.schema.v2/proto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

func main() {
	topic := "user-updates-value"
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

	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))
	if err != nil {
		log.Fatalf("Failed to create SR client: %s", err)
	}

	deser, err := protobuf.NewDeserializer(srClient, serde.ValueSerde, protobuf.NewDeserializerConfig())
	if err != nil {
		log.Fatalf("Failed to create deserializer: %s", err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			user := &pb.UserRequest{}
			err := deser.DeserializeInto(subject, msg.Value, user)
			if err != nil {
				log.Printf("Deserialization error: %s", err)
				continue
			}
			fmt.Printf("UserID: %#v, Name: %#v, Last name: %#v\n", user.UserId, user.Name, user.LastName)
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
