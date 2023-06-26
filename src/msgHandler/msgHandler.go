package msghandler

import (
	"encoding/binary"
	"strings"
)

type DyDanmuMsgHandler struct{}

func (h *DyDanmuMsgHandler) DyEncode(msg string) []byte {
	dataLen := len(msg) + 9
	lenByte := make([]byte, 4)
	msgByte := []byte(msg)
	sendByte := []byte{0xB1, 0x02, 0x00, 0x00}
	endByte := []byte{0x00}

	binary.LittleEndian.PutUint32(lenByte, uint32(dataLen))

	data := append(lenByte, lenByte...)
	data = append(data, sendByte...)
	data = append(data, msgByte...)
	data = append(data, endByte...)

	return data
}

func (h *DyDanmuMsgHandler) ParseMsg(rawMsg string) map[string]string {
	res := make(map[string]string)
	attrs := strings.Split(rawMsg, "/")[0 : len(strings.Split(rawMsg, "/"))-1]
	for _, attr := range attrs {
		attr = strings.ReplaceAll(attr, "@s", "/")
		attr = strings.ReplaceAll(attr, "@A", "@")
		couple := strings.Split(attr, "@=")
		res[couple[0]] = couple[1]
	}
	return res
}

func (h *DyDanmuMsgHandler) DyDecode(msgByte []byte) []string {
	pos := 0
	msg := make([]string, 0)
	for pos < len(msgByte) {
		contentLength := int(binary.LittleEndian.Uint32(msgByte[pos : pos+4]))
		content := string(msgByte[pos+12 : pos+3+contentLength])
		msg = append(msg, content)
		pos += 4 + contentLength
	}
	return msg
}

func (h *DyDanmuMsgHandler) GetChatMessages(msgByte []byte) []map[string]string {
	decodeMsg := h.DyDecode(msgByte)
	messages := make([]map[string]string, 0)
	for _, msg := range decodeMsg {
		res := h.ParseMsg(msg)
		if res["type"] != "chatmsg" {
			continue
		}
		messages = append(messages, res)
	}
	return messages
}
