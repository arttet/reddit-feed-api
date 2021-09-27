package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/arttet/reddit-feed-api/internal/config"
	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/opentracing/opentracing-go"

	"go.uber.org/zap"

	"github.com/Shopify/sarama"
)

type Producer interface {
	CreatePosts(offers []model.Post)
}

type producer struct {
	producer    sarama.SyncProducer
	topicName   string
	messageChan chan *sarama.ProducerMessage
	logger      *zap.Logger
}

type messageType uint16

const (
	created messageType = iota
)

type message struct {
	Type  messageType
	Value map[string]interface{}
}

func NewProducer(
	ctx context.Context,
	cfg *config.Kafka,
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

func (p *producer) CreatePosts(posts []model.Post) {
	result := make(map[string]interface{}, len(posts))
	for i, post := range posts {
		result[fmt.Sprintf("%d", i)] = post
	}

	if err := p.send("Producer.CreatePosts", created, result); err != nil {
		p.logger.Error("failed to send a message", zap.Error(err))
	}
}

func (p *producer) send(
	spanName string,
	msgType messageType,
	value map[string]interface{},
) error {

	span := opentracing.GlobalTracer().StartSpan(spanName)
	defer span.Finish()

	bytes, err := json.Marshal(
		message{
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
