package ports

import (
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
)

type DB interface {
	Migrate(models ...any) error

	CreateUser(user *entities.User) error
	GetUser(id uuid.UUID) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetAllUsers(offset, limit int) (*[]entities.User, error)
	UpdateUser(user *entities.User) error
	GetUserServers(userId uuid.UUID, offset, limit int) (*[]entities.ServerMember, error)
	GetUserDMChannels(userId uuid.UUID, offset, limit int) (*[]entities.DMChannelMember, error)
	GetUserChannelIds(userId uuid.UUID) (*[]uuid.UUID, error)
	GetUserServerIds(userId uuid.UUID) (*[]uuid.UUID, error) //TODO: implement this
	DeleteUser(id uuid.UUID) error

	CreateServer(server *entities.Server) error
	GetServer(id uuid.UUID) (*entities.Server, error)
	GetServerByName(name string) (*entities.Server, error)
	GetAllServers(offset, limit int) (*[]entities.Server, error)
	UpdateServer(server *entities.Server) error
	GetServerMembers(serverId uuid.UUID, offset, limit int) (*[]entities.ServerMember, error)
	AddServerMember(member *entities.ServerMember) error
	RemoveServerMember(serverId, userId uuid.UUID) error
	GetServerChannels(serverId uuid.UUID, offset, limit int) (*[]entities.ServerChannel, error)
	DeleteServer(id uuid.UUID) error

	CreateChannel(channel any) error
	GetChannel(channel any) error
	GetAllChannels(channels any, offset, limit int) error
	UpdateChannel(channel any) error
	GetChannelMembers(channelMembers any, channelId uuid.UUID, offset, limit int) error
	AddChannelMember(channelMember any) error
	RemoveChannelMember(channelMember any) error
	GetChannelMessages(channelMessages any, channelId uuid.UUID, offset, limit int) error
	DeleteChannel(channel any) error

	CreateMessage(msg any) error
	GetMessage(msg any) error
	UpdateMessage(msg any) error
	DeleteMessage(msg any) error
}
