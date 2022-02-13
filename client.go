package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = pongWait * 9 / 10
	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Client struct {
	conn *websocket.Conn
	hub  *Hub
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		if err := c.hub.rdb.Publish(context.Background(), "default", message).Err(); err != nil {
			log.Println(err)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	ctx := context.Background()
	sub := c.hub.rdb.Subscribe(ctx, "default")
	if _, err := sub.Receive(ctx); err != nil {
		log.Fatal(err)
	}
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
		_ = sub.Close()
	}()
	for {
		select {
		case msg := <-sub.Channel():
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write([]byte(msg.Payload)); err != nil {
				log.Println(err)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &Client{
		conn: conn,
		hub:  hub,
	}
	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}
