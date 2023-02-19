package dao

import (
	"errors"

	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"gorm.io/gorm"
)

type Message struct {
	db *gorm.DB
}

var (
	ErrNoSuchRecord = errors.New("no such video record")
)

// NewMessage create a message dao.
func NewMessage(db *gorm.DB) *Message {
	m := db.Migrator()
	if !m.HasTable(&model.Message{}) {
		err := m.CreateTable(&model.Message{})
		if err != nil {
			panic(err)
		}
	}
	return &Message{
		db: db,
	}
}

func (m *Message) GetMessages(toId, fromId, time int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := m.db.Model(model.Message{}).
		Where("to_user_id = ? AND from_user_id = ? AND create_time < ?", toId, fromId, time).
		Or("to_user_id = ? AND from_user_id = ? AND create_time < ?", fromId, toId, time).
		Order("create_time desc").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *Message) ChatAction(message *model.Message) error {
	err := m.db.Model(model.Message{}).
		Create(&message).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) GetLatestMessage(uId, toUId int64) (*model.Message, error) {
	var message *model.Message
	err := m.db.Model(model.Message{}).
		Where(&model.Message{ToUserId: uId, FromUserId: toUId}).Or(&model.Message{ToUserId: toUId, FromUserId: uId}).
		Order("create_time desc").First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoSuchRecord
		}
		return nil, err
	}
	return message, nil
}

func (m *Message) BatchGetLatestMessage(uId int64, toUIdList []int64) ([]*model.Message, error) {
	msgList := make([]*model.Message, 0)
	for _, toUid := range toUIdList {
		msg, err := m.GetLatestMessage(uId, toUid)
		if err != nil && err != ErrNoSuchRecord {
			return nil, err
		}
		msgList = append(msgList, msg)
	}
	return msgList, nil
}
