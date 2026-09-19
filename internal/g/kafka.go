// Copyright (C) edocseal. 2026-present.
//
// Created at 2026-09-19, by liasica

package g

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// kafkaTopic 应用日志统一投递的主题，由 ELK 的 Logstash 消费后按 logger 字段建索引
const kafkaTopic = "app-logs"

// kafkaWriter 把每条 JSON 日志异步投递到 Kafka，投递失败只写标准错误，不影响业务
type kafkaWriter struct {
	writer *kafka.Writer
}

// newKafkaWriter 创建异步投递器，主题不存在时由 broker 自动创建
func newKafkaWriter(brokers []string) *kafkaWriter {
	return &kafkaWriter{writer: &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  kafkaTopic,
		Balancer:               &kafka.LeastBytes{},
		Compression:            kafka.Snappy,
		Async:                  true,
		BatchTimeout:           time.Second,
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
		ErrorLogger: kafka.LoggerFunc(func(format string, args ...any) {
			_, _ = fmt.Fprintf(os.Stderr, "kafka 日志投递失败："+format+"\n", args...)
		}),
	}}
}

// Write 复制一条日志并异步发送，zap 会复用传入的缓冲区
func (w *kafkaWriter) Write(p []byte) (int, error) {
	message := bytes.Clone(bytes.TrimRight(p, "\n"))

	err := w.writer.WriteMessages(context.Background(), kafka.Message{Value: message})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "kafka 日志投递失败：%v\n", err)
	}

	return len(p), nil
}
