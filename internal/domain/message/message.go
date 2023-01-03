package message

import "time"

type Message struct {
	ID         uint      `json:"id" db:"id" gorm:"primary_key"`
	FromUserID uint      `json:"from_user_id" db:"from_user_id"`
	ToUserID   uint      `json:"to_user_id" db:"to_user_id"`
	Text       string    `json:"text" db:"text"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func (m Message) TableName() string {
	return "message"
}
