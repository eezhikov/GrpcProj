package kafka

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
	"time"
)

func StartKafka(chConn *sqlx.DB) {
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
		tx, err := chConn.Begin()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if _, err := tx.Exec("INSERT INTO users_log (info, action_time) VALUES ($1, $2)", string(m.Value), time.Now().Format("2006-01-02 15:04:05")); err != nil {
			fmt.Println(err)
			tx.Rollback()
			continue
		}
		if err := tx.Commit(); err != nil {
			fmt.Println(err)
			tx.Rollback()
			continue
		}
	}
}
