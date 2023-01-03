package message

import (
	"message/internal/domain/message"
)

type MessagesList struct {
	Items []message.Message `json:"items"`
	Count uint              `json:"count"`
}
