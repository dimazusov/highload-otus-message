package message

import (
	"context"
	"time"

	"message/internal/pkg/pagination"
)

type Service struct {
	rep *Repository
}

func New(rep *Repository) *Service {
	return &Service{
		rep: rep,
	}
}

func (m Service) GetMessage(ctx context.Context, id uint) (*Message, error) {
	msg, err := m.rep.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (m Service) GetMessages(ctx context.Context, cond *Message, pag *pagination.Pagination) ([]Message, error) {
	messages, err := m.rep.GetList(ctx, cond, pag)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m Service) SendMessage(ctx context.Context, fromUserID, toUserID uint, text string) error {
	_, err := m.rep.Create(ctx, &Message{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Text:       text,
		CreatedAt:  time.Now(),
	})
	return err
}

func (m Service) Update(ctx context.Context, msg *Message) error {
	return m.rep.Update(ctx, msg)
}

func (m Service) Delete(ctx context.Context, id uint) error {
	return m.rep.Delete(ctx, id)
}

func (m Service) Count(ctx context.Context, cond *Message) (uint, error) {
	l, err := m.rep.Count(ctx, cond)
	if err != nil {
		return 0, err
	}
	return l, nil
}
