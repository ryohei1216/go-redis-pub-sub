package graph

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/ryohei1216/go-redis-pub-sub/graph/model"
	"github.com/ryohei1216/go-redis-pub-sub/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	pubSubService service.PubSubService
	subscribers   map[string]chan<- *model.Message
	mutex         sync.Mutex
}

func NewResolver(pubSubService service.PubSubService) *Resolver {
	return &Resolver{
		pubSubService: pubSubService,
		subscribers:   map[string]chan<- *model.Message{},
		mutex:         sync.Mutex{},
	}
}

func (r *Resolver) SubscribeRedis(ctx context.Context) {
	go func() {
		pubsub := r.pubSubService.Subscribe(ctx, "messages")
		defer pubsub.Close()

		for msg := range pubsub.Channel() {
			// 受信したmessageはJSON形式なので、これをmodel.Message構造体に変換
			message := &model.Message{}
			err := json.Unmarshal([]byte(msg.Payload), message)
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			// 購読しているクライアントにRedisから受け取ったMessageをブロードキャスト
			r.mutex.Lock()
			for _, ch := range r.subscribers {
				ch <- message
			}
			r.mutex.Unlock()
		}
	}()

}
