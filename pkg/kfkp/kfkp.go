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

type KafkaConn struct {
	conn *kafka.Conn
}

type IKafkaConn interface {
	SendMessage(string) error
}

func NewConnection(ctx context.Context, cfg config.Configurations) (c IKafkaConn, err error) {
	conn, err := kafka.DialLeader(ctx, "tcp", cfg.KafkaURL, topic, partition)
	if err != nil {
		return
	}
	c = &KafkaConn{conn: conn}
	return
}

func (k *KafkaConn) SendMessage(message string) (err error) {
	_ = k.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = k.conn.WriteMessages(
		kafka.Message{Value: []byte(message)},
	)

	return
}
