package database

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
)

func (dbA *Adapter) CreateUser(user *entities.User) error {
	user.ID = uuid.New()

	return dbA.db.Create(user).Error
}

func (dbA *Adapter) GetUser(id uuid.UUID) (*entities.User, error) {
	user := &entities.User{ID: id}
	err := dbA.db.First(user).Error

	return user, err
}

func (dbA *Adapter) GetUserByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	err := dbA.db.First(user, "email = ?", email).Error

	return user, err
}

func (dbA *Adapter) GetAllUsers(offset, limit int) (*[]entities.User, error) {
	user := &[]entities.User{}
	err := dbA.db.Offset(offset).Limit(limit).Find(user).Error

	return user, err
}

func (dbA *Adapter) UpdateUser(user *entities.User) error {
	if user.ID == uuid.Nil {
		return errors.New("primary key must be specified")
	}

	return dbA.db.Model(user).Omit("ID", "CreatedAt").Updates(user).Error
}

func (dbA *Adapter) GetUserServers(userId uuid.UUID, offset, limit int) (*[]entities.Server, error) {
	serverMembers := &[]entities.ServerMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("server_id").
		Find(serverMembers, "user_id = ?", userId).Error

	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(*serverMembers))
	for idx, obj := range *serverMembers {
		ids[idx] = obj.ServerID
	}

	servers := &[]entities.Server{}
	err = dbA.db.Where("id IN ?", ids).Find(servers).Error

	return servers, err
}

func (dbA *Adapter) GetUserDMChannels(userId uuid.UUID, offset, limit int) (*[]entities.DMChannel, error) {
	channelMembers := &[]entities.DMChannelMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("channel_id").
		Find(channelMembers, "user_id = ?", userId).Error

	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(*channelMembers))
	for idx, obj := range *channelMembers {
		ids[idx] = obj.ChannelID
	}

	channels := &[]entities.DMChannel{}
	err = dbA.db.Where("id IN ?", ids).Find(channels).Error

	return channels, err
}

func (dbA *Adapter) GetUserChannelIds(userId uuid.UUID) (*[]uuid.UUID, error) {
	dmChannels := &[]entities.DMChannelMember{}
	err := dbA.db.Select("channel_id").
		Find(dmChannels, "user_id = ?", userId).Error

	serverChannels := &[]entities.ServerChannelMember{}
	err = dbA.db.Select("channel_id").
		Find(serverChannels, "user_id = ?", userId).Error

	ids := make([]uuid.UUID, len(*dmChannels)+len(*serverChannels))
	for idx, channel := range *dmChannels {
		ids[idx] = channel.ChannelID
	}

	offset := len(*dmChannels)
	for idx, channel := range *serverChannels {
		ids[offset+idx] = channel.ChannelID
	}

	return &ids, err
}

func (dbA *Adapter) GetUserServerIds(userId uuid.UUID) (*[]uuid.UUID, error) {
	servers := &[]entities.ServerMember{}
	err := dbA.db.Select("server_id").
		Find(servers, "user_id = ?", userId).Error

	ids := make([]uuid.UUID, len(*servers))
	for idx, server := range *servers {
		ids[idx] = server.ServerID
	}

	return &ids, err
}

func (dbA *Adapter) DeleteUser(userId uuid.UUID) error {
	user := &entities.User{ID: userId}

	return dbA.db.Delete(user).Error
}

func (dbA *Adapter) GetServerMemberRole(serverId, userId uuid.UUID) (string, error) {
	user := &entities.ServerMember{
		ServerID: serverId,
		UserID:   userId,
	}

	err := dbA.db.First(user).Error

	return user.Role, err
}
