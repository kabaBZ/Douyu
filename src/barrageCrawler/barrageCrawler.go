package barrageCrawler

import (
	"fmt"
	"time"

	"github.com/kabaBZ/Barrage_Go/src/barrageClient"
	msghandler "github.com/kabaBZ/Barrage_Go/src/msgHandler"
)

type DyDanmuCrawler struct {
	roomID         string
	heartbeatTimer *time.Timer
	client         *barrageClient.DyDanmuWebSocketClient
	msgHandler     *msghandler.DyDanmuMsgHandler
	keepHeartbeat  bool
}

func NewDyDanmuCrawler(roomID string) *DyDanmuCrawler {
	client := barrageClient.NewDyDanmuWebSocketClient()
	msgHandler := &msghandler.DyDanmuMsgHandler{}
	return &DyDanmuCrawler{
		roomID:         roomID,
		heartbeatTimer: nil,
		client:         client,
		msgHandler:     msgHandler,
		keepHeartbeat:  true,
	}

}

func (c *DyDanmuCrawler) Start() {
	c.client.Start()
	c.prepare()

	for {
		_, msg, _ := c.client.Websocket.ReadMessage()
		c.receiveMsg(msg)
	}

}

func (c *DyDanmuCrawler) Stop() {
	c.client.Stop()
	c.keepHeartbeat = false
}

func (c *DyDanmuCrawler) joinGroup() {
	joinGroupMsg := fmt.Sprintf("type@=joingroup/rid@=%s/gid@=1/", c.roomID)
	msgBytes := c.msgHandler.DyEncode(joinGroupMsg)
	c.client.Send(msgBytes)
}

func (c *DyDanmuCrawler) login() {
	loginMsg := fmt.Sprintf("type@=loginreq/roomid@=%s/dfl@=sn@AA=105@ASss@AA=1/username@=%s/uid@=%s/ver@=20190610/aver@=218101901/ct@=0/.",
		c.roomID, "99047358", "99047358")
	msgBytes := c.msgHandler.DyEncode(loginMsg)
	c.client.Send(msgBytes)
}

func (c *DyDanmuCrawler) startHeartbeat() {
	c.heartbeatTimer = time.NewTimer(0)
	go c.heartbeat()
}

func (c *DyDanmuCrawler) heartbeat() {
	heartbeatMsg := "type@=mrkl/"
	heartbeatMsgBytes := c.msgHandler.DyEncode(heartbeatMsg)
	for {
		select {
		case <-c.heartbeatTimer.C:
			c.client.Send(heartbeatMsgBytes)
			c.heartbeatTimer.Reset(30 * time.Second)
		default:
			time.Sleep(100 * time.Millisecond)
		}
		if !c.keepHeartbeat {
			return
		}
	}
}

func (c *DyDanmuCrawler) prepare() {
	c.login()
	c.joinGroup()
	c.startHeartbeat()
}

func (c *DyDanmuCrawler) receiveMsg(msg []byte) {
	chatMessages := c.msgHandler.GetChatMessages(msg)
	for _, message := range chatMessages {
		fmt.Printf("%s: %s\n", message["nn"], message["txt"])
	}
}
