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
	srClient, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://localhost:8081"))
	if err != nil {
		log.Fatalf("Failed to create SR client: %s", err)
	}

	serConfig := protobuf.NewSerializerConfig()
	ser, err := protobuf.NewSerializer(srClient, serde.ValueSerde, serConfig)
	if err != nil {
		log.Fatalf("Failed to create serializer: %s", err)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9093"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	topic := "user-updates-value"

	for i := 1; ; i++ {
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

		if i%1000 == 0 {
			p.Flush(15 * 10)
			fmt.Printf("Successfully produced %d messages.\n", i)
		}
	}

}
