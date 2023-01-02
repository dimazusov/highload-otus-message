package auth_token

import (
	"context"
	"message/internal/pkg/unsecure_jwt_token"

	"github.com/pkg/errors"
)

type service struct{}

const userIDKey = "userId"

var ErrWrongToken = errors.New("err wrong token")
var ErrorUserIDNotExists = errors.New("userId not exists")

type Service interface {
	Create(ctx context.Context, userID uint) (token string, err error)
	Parse(ctx context.Context, token string) (userID uint, err error)
}

func NewService() Service {
	return &service{}
}

func (m service) Create(ctx context.Context, userID uint) (token string, err error) {
	header := map[string]string{}
	payload := map[string]interface{}{userIDKey: userID}

	jwtToken, err := unsecure_jwt_token.Create(header, payload)
	if err != nil {
		return "", err
	}

	strJwtToken, err := jwtToken.ToString()
	if err != nil {
		return "", err
	}

	return strJwtToken, nil
}

func (m service) Parse(ctx context.Context, token string) (userID uint, err error) {
	jwtToken, err := unsecure_jwt_token.Parse(token)
	if err != nil {
		if errors.Is(err, unsecure_jwt_token.ErrWrongToken) {
			return 0, ErrWrongToken
		}
		return 0, err
	}

	userIDPayload, ok := jwtToken.GetPayload()[userIDKey]
	if !ok {
		return 0, ErrorUserIDNotExists
	}

	return uint(userIDPayload.(float64)), nil
}
