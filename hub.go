package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Hub struct {
	clients    map[*Client]bool
	rdb        *redis.Client
	register   chan *Client
	unregister chan *Client
}

func (h *Hub) run(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		case <-ctx.Done():
			return
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
