package database

import (
	"errors"
	"github.com/critch-app/critch-backend/internal/application/core/entities"
	"github.com/google/uuid"
)

func (dbA *Adapter) CreateServer(server *entities.Server) error {
	server.ID = uuid.New()

	return dbA.db.Create(server).Error
}

func (dbA *Adapter) GetServer(id uuid.UUID) (*entities.Server, error) {
	server := &entities.Server{ID: id}
	err := dbA.db.First(server).Error

	return server, err
}

func (dbA *Adapter) GetServerByName(name string) (*entities.Server, error) {
	server := &entities.Server{}
	err := dbA.db.First(server, "name = ?", name).Error

	return server, err
}

func (dbA *Adapter) GetAllServers(offset, limit int) (*[]entities.Server, error) {
	server := &[]entities.Server{}
	err := dbA.db.Offset(offset).Limit(limit).Find(server).Error

	return server, err
}

func (dbA *Adapter) UpdateServer(server *entities.Server) error {
	if server.ID == uuid.Nil {
		return errors.New("primary key must be specified")
	}

	return dbA.db.Model(server).Omit("ID", "CreatedAt").Updates(server).Error
}

func (dbA *Adapter) GetServerMembers(serverId uuid.UUID, offset, limit int) (*[]entities.User, error) {
	serverMembers := &[]entities.ServerMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("user_id", "role", "joined_at").
		Find(serverMembers, "server_id = ?", serverId).Error

	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(*serverMembers))
	for idx, obj := range *serverMembers {
		ids[idx] = obj.UserID
	}

	users := &[]entities.User{}
	err = dbA.db.Where("id IN ?", ids).Find(users).Error

	return users, err
}

func (dbA *Adapter) AddServerMember(member *entities.ServerMember) error {
	return dbA.db.Create(member).Error
}

func (dbA *Adapter) RemoveServerMember(serverId, userId uuid.UUID) error {
	member := &entities.ServerMember{
		ServerID: serverId,
		UserID:   userId,
	}

	return dbA.db.Delete(member).Error
}

func (dbA *Adapter) GetServerChannels(serverId, userId uuid.UUID, offset, limit int) (*[]entities.ServerChannel, error) {
	channelMember := &[]entities.ServerChannelMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("channel_id").
		Find(channelMember, "user_id = ? AND server_id = ?", userId, serverId).Error

	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(*channelMember))
	for idx, obj := range *channelMember {
		ids[idx] = obj.ChannelID
	}

	channels := &[]entities.ServerChannel{}
	err = dbA.db.Where("id IN ?", ids).Find(channels).Error

	return channels, err
}

func (dbA *Adapter) DeleteServer(id uuid.UUID) error {
	server := &entities.Server{ID: id}

	return dbA.db.Delete(server).Error
}
