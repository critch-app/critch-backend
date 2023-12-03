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

func (dbA *Adapter) GetUserServers(userId uuid.UUID, offset, limit int) (*[]entities.ServerMember, error) {
	servers := &[]entities.ServerMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("server_id", "role", "joined_at").
		Find(servers, "user_id = ?", userId).Error

	return servers, err
}

func (dbA *Adapter) GetUserDMChannels(userId uuid.UUID, offset, limit int) (*[]entities.DMChannelMember, error) {
	channels := &[]entities.DMChannelMember{}
	err := dbA.db.Offset(offset).Limit(limit).Select("channel_id").
		Find(channels, "user_id = ?", userId).Error

	return channels, err
}

func (dbA *Adapter) DeleteUser(userId uuid.UUID) error {
	user := &entities.User{ID: userId}

	return dbA.db.Delete(user).Error
}
