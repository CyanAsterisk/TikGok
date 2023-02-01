package dao

import (
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/global"
	"github.com/CyanAsterisk/TikGok/server/cmd/chat/model"
)

func GetMessages(toId int64, fromId int64) ([]*model.Message, error) {
	var messages []*model.Message
	err := global.DB.Model(model.Message{}).
		Where(&model.Message{ToUserId: toId, FromUserId: fromId}).Or(&model.Message{ToUserId: fromId, FromUserId: toId}).
		Order("create_date desc").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func ChatAction(message *model.Message) error {
	err := global.DB.Model(model.Message{}).
		Create(&message).Error
	if err != nil {
		return err
	}
	return nil
}

func GetLatestMessage(uId int64, toUId int64) (*model.Message, error) {
	var message *model.Message
	err := global.DB.Model(model.Message{}).
		Where(&model.Message{ToUserId: uId, FromUserId: toUId}).Or(&model.Message{ToUserId: toUId, FromUserId: uId}).
		Order("create_date desc").First(&message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}
