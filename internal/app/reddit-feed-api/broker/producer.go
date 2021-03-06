package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/Shopify/sarama"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

var (
	tracer trace.Tracer
)

func init() {
	tracer = otel.Tracer("github.com/arttet/reddit-feed-api/internal/app/reddit-feed-api/broker")
}

type Producer interface {
	CreatePosts(posts model.Posts)
}

type producer struct {
	producer    sarama.SyncProducer
	topicName   string
	messageChan chan *sarama.ProducerMessage
	logger      *zap.Logger
}

type MessageType uint16

const (
	Created MessageType = iota
)

type Message struct {
	Type  MessageType
	Value map[string]interface{}
}

func NewProducer(
	ctx context.Context,
	cfg *config.Producer,
	logger *zap.Logger,
) (
	Producer,
	error,
) {

	saramaCfg := sarama.NewConfig()

	saramaCfg.Producer.Partitioner = sarama.NewRandomPartitioner
	saramaCfg.Producer.RequiredAcks = sarama.WaitForAll
	saramaCfg.Producer.Return.Successes = true

	prod, err := sarama.NewSyncProducer(cfg.Brokers, saramaCfg)
	if err != nil {
		return nil, err
	}

	p := &producer{
		producer:    prod,
		topicName:   cfg.Topic,
		messageChan: make(chan *sarama.ProducerMessage, cfg.Capacity),
		logger:      logger,
	}

	go p.loop(ctx)

	return p, nil
}

func (p *producer) CreatePosts(posts model.Posts) {
	result := make(map[string]interface{}, len(posts))
	for i, post := range posts {
		result[fmt.Sprintf("%d", i)] = post
	}

	if err := p.send("Producer.CreatePosts", Created, result); err != nil {
		p.logger.Error("failed to send a message", zap.Error(err))
	}
}

func (p *producer) send(
	spanName string,
	msgType MessageType,
	value map[string]interface{},
) error {

	_, span := tracer.Start(context.Background(), spanName)
	defer span.End()

	bytes, err := json.Marshal(
		Message{
			Type:  msgType,
			Value: value,
		})

	if err != nil {
		return err
	}

	p.messageChan <- &sarama.ProducerMessage{
		Topic:     p.topicName,
		Key:       sarama.StringEncoder(p.topicName),
		Value:     sarama.StringEncoder(bytes),
		Partition: -1,
		Timestamp: time.Now(),
	}

	return nil
}

// loop fetches items using p.messageChan and sends them on SendMessage.
// loop exits when context is done.
func (p *producer) loop(ctx context.Context) {
	for {
		select {
		case msg := <-p.messageChan:
			partition, offset, err := p.producer.SendMessage(msg)

			if err != nil {
				p.logger.Error("failed to send a message", zap.Error(err))
				p.messageChan <- msg
			}

			p.logger.Info("delivered a message",
				zap.Int32("partition", partition),
				zap.String("topic", msg.Topic),
				zap.Int64("offset", offset),
			)

		case <-ctx.Done():
			close(p.messageChan)
			p.producer.Close()
			return
		}
	}
}
