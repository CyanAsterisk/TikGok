package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
	"gorm.io/gorm"
)

type Message struct {
	db *gorm.DB
}

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

func (m *Message) GetMessages(toId, fromId int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := m.db.Model(model.Message{}).
		Where(&model.Message{ToUserId: toId, FromUserId: fromId}).Or(&model.Message{ToUserId: fromId, FromUserId: toId}).
		Order("create_date desc").Find(&messages).Error
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
		Order("create_date desc").First(&message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *Message) BatchGetLatestMessage(uId int64, toUIdList []int64) ([]*model.Message, error) {
	msgList := make([]*model.Message, 0)
	for _, toUid := range toUIdList {
		msg, err := m.GetLatestMessage(uId, toUid)
		if err != nil {
			return nil, err
		}
		msgList = append(msgList, msg)
	}
	return msgList, nil
}
