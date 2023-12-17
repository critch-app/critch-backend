package database

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
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

func (dbA *Adapter) GetServerMembers(serverId uuid.UUID, offset, limit int) (*[]entities.ServerMember, error) {
	members := &[]entities.ServerMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("user_id", "role", "joined_at").
		Find(members, "server_id = ?", serverId).Error

	return members, err
}

func (dbA *Adapter) AddServerMember(serverId, userId uuid.UUID) error {
	member := &entities.ServerMember{
		ServerID: serverId,
		UserID:   userId,
	}

	return dbA.db.Create(member).Error
}

func (dbA *Adapter) RemoveServerMember(serverId, userId uuid.UUID) error {
	member := &entities.ServerMember{
		ServerID: serverId,
		UserID:   userId,
	}

	return dbA.db.Delete(member).Error
}

func (dbA *Adapter) GetServerChannels(serverId uuid.UUID, offset, limit int) (*[]entities.ServerChannel, error) {
	channels := &[]entities.ServerChannel{}
	err := dbA.db.Offset(offset).Limit(limit).Omit("server_id").
		Find(channels, "server_id = ?", serverId).Error

	return channels, err
}

func (dbA *Adapter) DeleteServer(id uuid.UUID) error {
	server := &entities.Server{ID: id}

	return dbA.db.Delete(server).Error
}
