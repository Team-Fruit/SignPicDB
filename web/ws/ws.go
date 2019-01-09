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

const pingPeriod = 54 * time.Second

var (
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
		case <-ticker.C:
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
