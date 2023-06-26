package barrageClient

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type DyDanmuWebSocketClient struct {
	url       string
	Websocket *websocket.Conn
}

func NewDyDanmuWebSocketClient() *DyDanmuWebSocketClient {
	url := "wss://danmuproxy.douyu.com:8506/"
	return &DyDanmuWebSocketClient{
		url:       url,
		Websocket: nil,
	}
}

func (c *DyDanmuWebSocketClient) Start() {
	u, err := url.Parse(c.url)
	if err != nil {
		log.Fatal(err)
	}

	header := http.Header{}
	header.Add("Origin", "https://www.douyu.com")

	websocket.DefaultDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal(err)
	}
	c.Websocket = conn

}

func (c *DyDanmuWebSocketClient) Stop() {
	c.Websocket.Close()
}

func (c *DyDanmuWebSocketClient) Send(msg []byte) {
	c.Websocket.WriteMessage(websocket.BinaryMessage, msg)
}
