package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func StartKafka() {

	conf := kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "idroot",
		GroupID:  "g1",
		MaxBytes: 10,
	}
	reader := kafka.NewReader(conf)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println("message: ", string(m.Value))
	}
}
