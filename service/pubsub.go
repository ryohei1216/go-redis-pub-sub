package service

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ErrPublish = errors.New("failed to publish")

type PubSubService interface {
	// Publish publish message to channel.
	Publish(ctx context.Context, channel string, payload any) error

	// Subscribe return PubSub struct that has a method called ReceiveMessage().
	Subscribe(ctx context.Context, channel string) *redis.PubSub
}

type pubSubService struct {
	client *redis.Client
}

func NewPubSubService(client *redis.Client) *pubSubService {
	return &pubSubService{
		client: client,
	}
}

// Publish publish message to channel.
func (s *pubSubService) Publish(ctx context.Context, channel string, payload any) error {
	err := s.client.Publish(ctx, channel, payload).Err()

	if err != nil {
		return ErrPublish
	}

	return nil
}

// Subscribe return PubSub struct that has a method called ReceiveMessage().
func (s *pubSubService) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	pubSub := s.client.Subscribe(ctx, channel)

	defer pubSub.Close()

	return pubSub
}
