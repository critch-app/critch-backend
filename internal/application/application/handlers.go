package application

import (
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
)

func (app App) Login(username, password string) error {
	//TODO implement me
	return nil
}

func (app App) Signup(user *entities.User) error {
	//TODO hash password first
	return app.db.CreateUser(user)
}

func (app App) GetUser(id uuid.UUID) (*entities.User, error) {
	return app.db.GetUser(id)
}

func (app App) GetUserByEmail(email string) (*entities.User, error) {
	return app.db.GetUserByEmail(email)
}

func (app App) GetAllUsers(offset, limit int) (*[]entities.User, error) {
	return app.db.GetAllUsers(offset, limit)
}

func (app App) UpdateUser(user *entities.User) error {
	return app.db.UpdateUser(user)
}

func (app App) GetUserServers(userId uuid.UUID, offset, limit int) (*[]entities.ServerMember, error) {
	return app.db.GetUserServers(userId, offset, limit)
}

func (app App) GetUserDMChannels(userId uuid.UUID, offset, limit int) (*[]entities.DMChannelMember, error) {
	return app.db.GetUserDMChannels(userId, offset, limit)
}

func (app App) DeleteUser(id uuid.UUID) error {
	return app.db.DeleteUser(id)
}

func (app App) CreateServer(server *entities.Server) error {
	return app.db.CreateServer(server)
}

func (app App) GetServer(id uuid.UUID) (*entities.Server, error) {
	return app.db.GetServer(id)
}

func (app App) GetServerByName(name string) (*entities.Server, error) {
	return app.db.GetServerByName(name)
}

func (app App) GetAllServers(offset, limit int) (*[]entities.Server, error) {
	return app.db.GetAllServers(offset, limit)
}

func (app App) UpdateServer(server *entities.Server) error {
	return app.db.UpdateServer(server)
}

func (app App) GetServerMembers(serverId uuid.UUID, offset, limit int) (*[]entities.ServerMember, error) {
	return app.db.GetServerMembers(serverId, offset, limit)
}

func (app App) AddServerMember(serverId, userId uuid.UUID) error {
	return app.db.AddServerMember(serverId, userId)
}

func (app App) RemoveServerMember(serverId, userId uuid.UUID) error {
	return app.db.RemoveServerMember(serverId, userId)
}

func (app App) GetServerChannels(serverId uuid.UUID, offset, limit int) (*[]entities.ServerChannel, error) {
	return app.db.GetServerChannels(serverId, offset, limit)
}

func (app App) DeleteServer(id uuid.UUID) error {
	return app.db.DeleteServer(id)
}

func (app App) CreateChannel(channel any) error {
	return app.db.CreateChannel(channel)
}

func (app App) GetChannel(channel any) error {
	return app.db.GetChannel(channel)
}

func (app App) GetAllChannels(channels any, offset, limit int) error {
	return app.db.GetAllChannels(channels, offset, limit)
}

func (app App) UpdateChannel(channel any) error {
	return app.db.UpdateChannel(channel)
}

func (app App) GetChannelMembers(channelMembers any, channelId uuid.UUID, offset, limit int) error {
	return app.db.GetChannelMembers(channelMembers, channelId, offset, limit)
}

func (app App) AddChannelMember(channelMember any) error {
	return app.db.AddChannelMember(channelMember)
}

func (app App) RemoveChannelMember(channelMember any) error {
	return app.db.RemoveChannelMember(channelMember)
}

func (app App) GetChannelMessages(channelMessages any, channelId uuid.UUID, offset, limit int) error {
	return app.db.GetChannelMessages(channelMessages, channelId, offset, limit)
}

func (app App) DeleteChannel(channel any) error {
	return app.db.DeleteChannel(channel)
}

func (app App) GetMessage(msg any) error {
	return app.db.GetMessage(msg)
}

func (app App) UpdateMessage(msg any) error {
	return app.db.UpdateMessage(msg)
}

func (app App) DeleteMessage(msg any) error {
	return app.db.DeleteMessage(msg)
}
