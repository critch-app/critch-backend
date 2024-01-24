package database

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
)

func (dbA *Adapter) CreateMessage(msg any) error {
	err := addMessageID(msg)
	if err != nil {
		return err
	}

	return dbA.db.Create(msg).Error
}

func (dbA *Adapter) GetMessage(msg any) error {
	err := checkMessageID(msg)
	if err != nil {
		return err
	}

	return dbA.db.First(msg).Error
}

func (dbA *Adapter) UpdateMessage(msg any) error {
	err := checkMessageID(msg)
	if err != nil {
		return err
	}

	return dbA.db.Model(msg).Select("content", "attachment").Updates(msg).Error
}

func (dbA *Adapter) DeleteMessage(msg any) error {
	err := checkMessageID(msg)
	if err != nil {
		return err
	}

	return dbA.db.Delete(msg).Error
}

func addMessageID(msg any) error {
	err := validateMessageType(msg)
	if err != nil {
		return err
	}

	switch msg.(type) {
	case *entities.ServerMessage:
		messageModel := msg.(*entities.ServerMessage)
		messageModel.ID = uuid.New()
	case *entities.DirectMessage:
		messageModel := msg.(*entities.DirectMessage)
		messageModel.ID = uuid.New()
	}

	return nil
}

func validateMessageType(message any) error {
	switch message.(type) {
	case *entities.ServerMessage:
	case *entities.DirectMessage:
	default:
		return errors.New("message must be of type *ServerMessage or *DirectMessage")
	}

	return nil
}

func checkMessageID(message any) error {
	err := validateMessageType(message)
	if err != nil {
		return err
	}

	switch message.(type) {
	case *entities.ServerMessage:
		messageModel := message.(*entities.ServerMessage)
		if messageModel.ID == uuid.Nil {
			return errors.New("primary key must be specified")
		}
	case *entities.DirectMessage:
		messageModel := message.(*entities.DirectMessage)
		if messageModel.ID == uuid.Nil {
			return errors.New("primary key must be specified")
		}
	}

	return nil
}
