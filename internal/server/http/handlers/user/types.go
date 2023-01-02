package user

import (
	"message/internal/domain/user"
)

type UsersList struct {
	Items []user.User `json:"items"`
	Count uint        `json:"count"`
}
