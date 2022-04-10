package broker

import (
	"context"
	"encoding/json"

	"github.com/arttet/reddit-feed-api/internal/consume"
	"github.com/arttet/reddit-feed-api/internal/model"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"

	"go.uber.org/zap"
)

func ReceiveMessage(logger *zap.Logger) consume.ConsumeFunction {
	return func(ctx context.Context, message *sarama.ConsumerMessage) error {
		var msg Message

		if err := json.Unmarshal(message.Value, &msg); err != nil {
			logger.Error("failed to unmarshal a message", zap.Error(err))

			return err
		}

		for i := range msg.Value {
			var post model.Post
			if err := mapstructure.Decode(msg.Value[i], &post); err != nil {
				logger.Error("failed to unmarshal a message", zap.Error(err))

				return err
			}

			logger.Info("the message has been received",
				zap.Uint16("Type", uint16(msg.Type)),
				zap.Reflect("Posts", post),
			)
		}

		return nil
	}
}
