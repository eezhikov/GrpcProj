package logger

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type KafkaProducerConf struct {
	namespace, broker, topic string
	debug, async             bool
}

var prodCfg = KafkaProducerConf{
	"test-service",
	"localhost:9092",
	"idroot",
	false,
	true,
}

type KafkaProducer struct {
	sync.Mutex

	Producer *kafka.Writer
}

var kp = &KafkaProducer{
	Producer: &kafka.Writer{
		Addr:         kafka.TCP(prodCfg.broker),
		Topic:        prodCfg.topic,
		Async:        prodCfg.async,
		BatchTimeout: 1,
		RequiredAcks: kafka.RequireAll,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}

		},
	},
}

var customConfig = zapcore.EncoderConfig{
	TimeKey:        "timeStamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
}

func NewLogger() *zap.Logger {
	return getLogger(kp)
}

func getLogger(kp *KafkaProducer) *zap.Logger {
	// cores are logger interfaces
	var core zapcore.Core

	core = zapcore.NewCore(zapcore.NewJSONEncoder(customConfig), zapcore.Lock(zapcore.AddSync(kp)), zap.DebugLevel)

	// join inputs, encoders, level-handling functions into cores, then "tee" together
	logger := zap.New(core)
	defer logger.Sync()
	return logger
}

func (kp *KafkaProducer) Write(msg []byte) (int, error) {
	kp.Lock()
	defer kp.Unlock()
	err := kp.Producer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(""),
		Value: msg,
	})

	return len(msg), err
}
