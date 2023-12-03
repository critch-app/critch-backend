package database

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
)

func (dbA *Adapter) CreateChannel(channel any) error {
	err := addChannelID(channel)
	if err != nil {
		return err
	}

	return dbA.db.Create(channel).Error
}

func (dbA *Adapter) GetChannel(channel any) error {
	err := checkChannelID(channel)
	if err != nil {
		return err
	}

	return dbA.db.First(channel).Error
}

func (dbA *Adapter) GetAllChannels(channels any, offset, limit int) error {
	err := validateChannelArrayType(channels)
	if err != nil {
		return err
	}

	return dbA.db.Offset(offset).Limit(limit).Find(channels).Error
}

func (dbA *Adapter) UpdateChannel(channel any) error {
	err := checkChannelID(channel)
	if err != nil {
		return err
	}

	return dbA.db.Model(channel).Select("name", "description").Updates(channel).Error
}

func (dbA *Adapter) GetChannelMembers(channelMembers any, channelId uuid.UUID, offset, limit int) error {
	err := validateChannelMembersType(channelMembers)
	if err != nil {
		return err
	}

	return dbA.db.Offset(offset).Limit(limit).Select("user_id").
		Find(channelMembers, "channel_id = ?", channelId).Error
}

func (dbA *Adapter) AddChannelMember(channelMember any) error {
	err := validateChannelMembersType(channelMember)
	if err != nil {
		return err
	}

	return dbA.db.Create(channelMember).Error
}

func (dbA *Adapter) RemoveChannelMember(channelMember any) error {
	err := validateChannelMembersType(channelMember)
	if err != nil {
		return err
	}

	return dbA.db.Delete(channelMember).Error
}

func (dbA *Adapter) GetChannelMessages(channelMessages any, offset, limit int) error {
	err := validateChannelMessageType(channelMessages)
	if err != nil {
		return err
	}

	return dbA.db.Offset(offset).Limit(limit).Order("created_at").Find(channelMessages).Error
}

func (dbA *Adapter) DeleteChannel(channel any) error {
	err := checkChannelID(channel)
	if err != nil {
		return err
	}

	return dbA.db.Delete(channel).Error
}

func validateChannelType(channel any) error {
	switch channel.(type) {
	case *entities.ServerChannel:
	case *entities.DMChannel:
	default:
		return errors.New("channel must be of type *ServerChannel or *DMChannel")
	}

	return nil
}

func validateChannelArrayType(channels any) error {
	switch channels.(type) {
	case *[]entities.ServerChannel:
	case *[]entities.DMChannel:
	default:
		return errors.New("channels must be of type *[]ServerChannel or *[]DMChannel")
	}

	return nil
}

func validateChannelMembersType(members any) error {
	switch members.(type) {
	case *[]entities.ServerChannelMember:
	case *[]entities.DMChannelMember:
	default:
		return errors.New("channels must be of type *[]ServerChannelMember or *[]DMChannelMember")
	}

	return nil
}

func validateChannelMessageType(messages any) error {
	switch messages.(type) {
	case *[]entities.ServerMessage:
	case *[]entities.DirectMessage:
	default:
		return errors.New("channels must be of type *[]ServerChannelMember or *[]DMChannelMember")
	}

	return nil
}

func addChannelID(channel any) error {
	err := validateChannelType(channel)
	if err != nil {
		return err
	}

	switch channel.(type) {
	case *entities.ServerChannel:
		channelModel := channel.(*entities.ServerChannel)
		channelModel.ID = uuid.New()
	case *entities.DMChannel:
		channelModel := channel.(*entities.DMChannel)
		channelModel.ID = uuid.New()
	}

	return nil
}

func checkChannelID(channel any) error {
	err := validateChannelType(channel)
	if err != nil {
		return err
	}

	switch channel.(type) {
	case *entities.ServerChannel:
		channelModel := channel.(*entities.ServerChannel)
		if channelModel.ID == uuid.Nil {
			return errors.New("primary key must be specified")
		}
	case *entities.DMChannel:
		channelModel := channel.(*entities.DMChannel)
		if channelModel.ID == uuid.Nil {
			return errors.New("primary key must be specified")
		}
	}

	return nil
}
