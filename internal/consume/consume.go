package consume

import (
	"context"
	"time"

	"github.com/arttet/reddit-feed-api/internal/config"

	"go.uber.org/zap"

	"github.com/Shopify/sarama"
)

type ConsumeFunction func(ctx context.Context, message *sarama.ConsumerMessage) error

type Consumer interface {
	sarama.ConsumerGroupHandler
}

type consumer struct {
	fn ConsumeFunction
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		func(message *sarama.ConsumerMessage) {
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			defer cancel()

			if err := c.fn(ctx, message); err == nil {
				session.MarkMessage(message, "")
			}
		}(message)
	}

	return nil
}

func StartConsuming(
	ctx context.Context,
	cfg config.Consume,
	consumeFunction ConsumeFunction,
	logger *zap.Logger,
) error {

	saramaCfg := sarama.NewConfig()

	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, saramaCfg)
	if err != nil {
		return err
	}

	consumer := &consumer{
		fn: consumeFunction,
	}

	for {
		if err := consumerGroup.Consume(ctx, []string{cfg.Topic}, consumer); err != nil {
			logger.Error("failed to consume a message", zap.Error(err))
		}

		if ctx.Err() != nil {
			break
		}
	}

	if err = consumerGroup.Close(); err != nil {
		logger.Error("failed to close a consumer group", zap.Error(err))
	}

	return nil
}
