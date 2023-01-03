package message

import (
	"context"

	"gorm.io/gorm"

	"message/internal/pkg/pagination"
)

const createQuery = "INSERT INTO message VALUES (nextval('serial_message_id'), ?,?,?,?,null)"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (m Repository) Get(ctx context.Context, messageID uint) (*Message, error) {
	msg := &Message{}
	err := m.db.WithContext(ctx).First(msg, messageID).Error
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (m Repository) GetList(ctx context.Context, cond *Message, pag *pagination.Pagination) ([]Message, error) {
	messages := []Message{}
	pag.GetLimit()
	pag.GetOffset()
	err := m.db.WithContext(ctx).
		Where(cond).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m Repository) Create(ctx context.Context, msg *Message) (uint, error) {
	err := m.db.WithContext(ctx).Exec(createQuery,
		msg.FromUserID,
		msg.ToUserID,
		msg.Text,
		msg.CreatedAt).Error
	if err != nil {
		return 0, err
	}
	return msg.ID, nil
}

func (m Repository) Update(ctx context.Context, msg *Message) error {
	err := m.db.WithContext(ctx).Save(msg).Error
	if err != nil {
		return nil
	}
	return nil
}

func (m Repository) Delete(ctx context.Context, id uint) error {
	return m.db.WithContext(ctx).Delete(&Message{ID: id}).Error
}

func (m Repository) Count(ctx context.Context, cond *Message) (uint, error) {
	var count int64

	err := m.db.WithContext(ctx).Where(cond).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}
