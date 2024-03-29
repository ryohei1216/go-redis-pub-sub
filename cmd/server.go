package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/ryohei1216/go-redis-pub-sub/graph"
	"github.com/ryohei1216/go-redis-pub-sub/redis"
	"github.com/ryohei1216/go-redis-pub-sub/service"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx := context.Background()
	redisClient := redis.New(ctx)

	pubSubService := service.NewPubSubService(redisClient)

	resolver := graph.NewResolver(pubSubService)
	resolver.SubscribeRedis(ctx)

	srv := handler.New(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}),
	)

	// Subscription を使うには `transport.POST` と `transport.Websocket` をトランスポートとして追加すればよい
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
