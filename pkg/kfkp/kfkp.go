package kfkp

import (
	"context"
	"time"

	"github.com/danielboakye/go-xm/config"
	"github.com/segmentio/kafka-go"
)

const (
	topic     = "companies"
	partition = 0
)

func SendMessage(ctx context.Context, cfg config.Configurations, message string) (err error) {

	var conn *kafka.Conn
	conn, err = kafka.DialLeader(ctx, "tcp", cfg.KafkaURL, topic, partition)
	if err != nil {
		return
	}

	_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)

	if err != nil {
		return
	}

	err = conn.Close()

	return
}
