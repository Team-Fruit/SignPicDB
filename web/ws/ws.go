package ws

import (
	"net/http"
	"time"
	"encoding/json"

	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"github.com/Team-Fruit/SignPicDB/web/models"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

const (
	pingPeriod = 54 * time.Second
	writeWait = 10 * time.Second
)

var (
	newline = []byte{'\n'}

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	AnalyticsChan = make(chan models.AnalyticsData)
)

func (c *Client) pingTicker() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) broadcast() {
	for {
		select {
		case p := <-AnalyticsChan:
			b, _ := json.Marshal(p)
			c.hub.broadcast <- b
		}
	}
}

func ServeWs(hub *Hub, c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.pingTicker()
	go client.broadcast()

	return nil
}
