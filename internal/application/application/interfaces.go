package application

import (
	"github.com/critch-app/critch-backend/internal/application/core/entities"
	"github.com/critch-app/critch-backend/internal/application/core/msgsrvc"
	"github.com/google/uuid"
)

type AppI interface {
	Login(email, password string) (string, uuid.UUID, error)
	Signup(user *entities.User) error
	GetUser(id uuid.UUID) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetAllUsers(offset, limit int) (*[]entities.User, error)
	UpdateUser(user *entities.User) error
	GetUserServers(userId uuid.UUID, offset, limit int) (*[]entities.Server, error)
	GetUserDMChannels(userId uuid.UUID, offset, limit int) (*[]entities.DMChannel, error)
	DeleteUser(id uuid.UUID) error

	CreateServer(server *entities.Server, OwnerID uuid.UUID) error
	GetServer(id uuid.UUID) (*entities.Server, error)
	GetServerByName(name string) (*entities.Server, error)
	GetAllServers(offset, limit int) (*[]entities.Server, error)
	UpdateServer(server *entities.Server) error
	GetServerMembers(serverId uuid.UUID, offset, limit int) (*[]entities.User, error)
	AddServerMember(serverId, userId uuid.UUID) error
	RemoveServerMember(serverId, userId uuid.UUID) error
	GetServerChannels(serverId, userId uuid.UUID, offset, limit int) (*[]entities.ServerChannel, error)
	DeleteServer(id uuid.UUID) error

	CreateChannel(channel any, userId uuid.UUID, isServerChannel bool) error
	GetChannel(channel any) error
	GetAllChannels(channels any, offset, limit int) error
	UpdateChannel(channel any) error
	GetChannelMembers(channelMembers any, channelId uuid.UUID, offset, limit int) error
	AddChannelMember(channelMember any) error
	RemoveChannelMember(channelMember any) error
	GetChannelMessages(channelMessages any, channelId uuid.UUID, offset, limit int) error
	DeleteChannel(channel any) error

	GetMessage(msg any) error
	UpdateMessage(msg any) error
	DeleteMessage(msg any) error

	ValidateJWTToken(tokenString string) (uuid.UUID, error)

	SendMessages(incomingMessage *msgsrvc.IncomingMessage) error
	SendNotification(notificationObj any, serverId uuid.UUID) error
	ReceiveMessages(client *msgsrvc.Client) (any, bool)

	ConnectWebsocket(clientId uuid.UUID) (*msgsrvc.Client, error)
	JoinChannels(clientObj *msgsrvc.Client, serverId uuid.UUID, channels []uuid.UUID)
	QuitChannel(clientObj *msgsrvc.Client, channelId uuid.UUID)
	QuitServer(clientObj *msgsrvc.Client, serverId uuid.UUID)
	RemoveChannel(channelId uuid.UUID)
	RemoveServer(serverId uuid.UUID)
	DisconnectWebsocket(client *msgsrvc.Client)

	GetServerMemberRole(serverId, userId uuid.UUID) (string, error)
}
