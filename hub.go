package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	rdb        *redis.Client
	register   chan *Client
	unregister chan *Client
}

func (h *Hub) run() {
	ctx := context.Background()
	sub := h.rdb.Subscribe(ctx, "default")
	defer func() {
		_ = sub.Close()
	}()
	if _, err := sub.Receive(ctx); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.messages)
			}
		case msg := <-sub.Channel():
			for client := range h.clients {
				select {
				case client.messages <- []byte(msg.Payload):
				default:
					delete(h.clients, client)
					close(client.messages)
				}
			}
		}
	}
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rdb:        NewRDB(),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
