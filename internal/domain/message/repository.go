package message

import "database/sql"

type Repository struct{}

func NewRepository([]sql.DB) *Repository {
	return &Repository{}
}

func GetMessages(userID int) ([]Message, error) {
	return nil, nil
}

func SendMessage(fromUserID, toUserID int, text string) error {
	return nil
}
