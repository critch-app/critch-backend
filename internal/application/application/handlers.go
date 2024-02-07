package application

import (
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/msgsrvc"
	"golang.org/x/crypto/bcrypt"
)

func (app *App) Login(email, password string) (string, uuid.UUID, error) {
	user, err := app.db.GetUserByEmail(email)
	if err != nil {
		return "", uuid.Nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", uuid.Nil, err
	}

	token, err := newJWTToken(user.ID)
	if err != nil {
		return "", uuid.Nil, err
	}

	return token, user.ID, nil
}

func (app *App) Signup(user *entities.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return app.db.CreateUser(user)
}

func (app *App) GetUser(id uuid.UUID) (*entities.User, error) {
	return app.db.GetUser(id)
}

func (app *App) GetUserByEmail(email string) (*entities.User, error) {
	return app.db.GetUserByEmail(email)
}

func (app *App) GetAllUsers(offset, limit int) (*[]entities.User, error) {
	return app.db.GetAllUsers(offset, limit)
}

func (app *App) UpdateUser(user *entities.User) error {
	return app.db.UpdateUser(user)
}

func (app *App) GetUserServers(userId uuid.UUID, offset, limit int) (*[]entities.Server, error) {
	return app.db.GetUserServers(userId, offset, limit)
}

func (app *App) GetUserDMChannels(userId uuid.UUID, offset, limit int) (*[]entities.DMChannel, error) {
	return app.db.GetUserDMChannels(userId, offset, limit)
}

func (app *App) DeleteUser(id uuid.UUID) error {
	return app.db.DeleteUser(id)
}

func (app *App) CreateServer(server *entities.Server, OwnerID uuid.UUID) error {
	err := app.db.CreateServer(server)
	if err != nil {
		return err
	}

	err = app.db.AddServerMember(&entities.ServerMember{
		ServerID: server.ID,
		UserID:   OwnerID,
		Role:     "owner",
	})

	return err
}

func (app *App) GetServer(id uuid.UUID) (*entities.Server, error) {
	return app.db.GetServer(id)
}

func (app *App) GetServerByName(name string) (*entities.Server, error) {
	return app.db.GetServerByName(name)
}

func (app *App) GetAllServers(offset, limit int) (*[]entities.Server, error) {
	return app.db.GetAllServers(offset, limit)
}

func (app *App) UpdateServer(server *entities.Server) error {
	return app.db.UpdateServer(server)
}

func (app *App) GetServerMembers(serverId uuid.UUID, offset, limit int) (*[]entities.User, error) {
	return app.db.GetServerMembers(serverId, offset, limit)
}

func (app *App) AddServerMember(serverId, userId uuid.UUID) error {
	return app.db.AddServerMember(&entities.ServerMember{
		ServerID: serverId,
		UserID:   userId,
		Role:     "member",
	})
}

func (app *App) RemoveServerMember(serverId, userId uuid.UUID) error {
	return app.db.RemoveServerMember(serverId, userId)
}

func (app *App) GetServerChannels(serverId, userId uuid.UUID, offset, limit int) (*[]entities.ServerChannel, error) {
	return app.db.GetServerChannels(serverId, userId, offset, limit)
}

func (app *App) DeleteServer(id uuid.UUID) error {
	return app.db.DeleteServer(id)
}

func (app *App) CreateChannel(channel any, userId uuid.UUID, isServerChannel bool) error {
	err := app.db.CreateChannel(channel)
	if err != nil {
		return err
	}

	var channelMember any
	if isServerChannel {
		channelMember = &entities.ServerChannelMember{
			ChannelID: channel.(*entities.ServerChannel).Channel.ID,
			ServerID:  channel.(*entities.ServerChannel).ServerID,
			UserID:    userId,
		}
	} else {
		channelMember = &entities.DMChannelMember{
			ChannelID: channel.(*entities.DMChannel).Channel.ID,
			UserID:    userId,
		}
	}

	err = app.db.AddChannelMember(channelMember)

	return err
}

func (app *App) GetChannel(channel any) error {
	return app.db.GetChannel(channel)
}

func (app *App) GetAllChannels(channels any, offset, limit int) error {
	return app.db.GetAllChannels(channels, offset, limit)
}

func (app *App) UpdateChannel(channel any) error {
	return app.db.UpdateChannel(channel)
}

func (app *App) GetChannelMembers(channelMembers any, channelId uuid.UUID, offset, limit int) error {
	return app.db.GetChannelMembers(channelMembers, channelId, offset, limit)
}

func (app *App) AddChannelMember(channelMember any) error {
	return app.db.AddChannelMember(channelMember)
}

func (app *App) RemoveChannelMember(channelMember any) error {
	return app.db.RemoveChannelMember(channelMember)
}

func (app *App) GetChannelMessages(channelMessages any, channelId uuid.UUID, offset, limit int) error {
	return app.db.GetChannelMessages(channelMessages, channelId, offset, limit)
}

func (app *App) DeleteChannel(channel any) error {
	return app.db.DeleteChannel(channel)
}

func (app *App) GetMessage(msg any) error {
	return app.db.GetMessage(msg)
}

func (app *App) UpdateMessage(msg any) error {
	return app.db.UpdateMessage(msg)
}

func (app *App) DeleteMessage(msg any) error {
	return app.db.DeleteMessage(msg)
}

func (app *App) SendMessages(incomingMessage *msgsrvc.IncomingMessage) error {
	var outgoingMessage any
	if incomingMessage.ServerId != uuid.Nil {
		outgoingMessage = &entities.ServerMessage{
			Message: entities.Message{
				ChannelID:  incomingMessage.ChannelId,
				SenderID:   incomingMessage.SenderId,
				Content:    incomingMessage.Content,
				Attachment: incomingMessage.Attachment,
			},
		}
	} else {
		outgoingMessage = &entities.DirectMessage{
			Message: entities.Message{
				ChannelID:  incomingMessage.ChannelId,
				SenderID:   incomingMessage.SenderId,
				Content:    incomingMessage.Content,
				Attachment: incomingMessage.Attachment,
			},
		}
	}

	err := app.db.CreateMessage(outgoingMessage)
	if err != nil {
		return err
	}

	app.messagingService.Broadcast <- &msgsrvc.BroadcastMessage{
		Type:      msgsrvc.MESSAGE,
		ChannelId: incomingMessage.ChannelId,
		ServerId:  incomingMessage.ServerId,
		Message:   outgoingMessage,
	}

	return nil
}

func (app *App) SendNotification(notificationObj any) error {
	app.messagingService.Broadcast <- &msgsrvc.BroadcastMessage{
		Type:    msgsrvc.NOTIFICATION,
		Message: notificationObj,
	}

	return nil
}

func (app *App) ReceiveMessages(client *msgsrvc.Client) (any, bool) {
	msg, ok := <-client.MessagingChannel
	return msg, ok
}

func (app *App) ConnectWebsocket(clientId uuid.UUID) (*msgsrvc.Client, error) {
	channelIds, err := app.db.GetUserChannelIds(clientId)
	if err != nil {
		return nil, err
	}

	serverIds, err := app.db.GetUserServerIds(clientId)
	if err != nil {
		return nil, err
	}

	client := &msgsrvc.Client{
		ID:               clientId,
		MessagingChannel: make(chan any, 10),
	}

	app.messagingService.Connect <- &msgsrvc.NewClient{
		ClientObj: client,
		Servers:   serverIds,
		Channels:  channelIds,
	}

	return client, nil
}

func (app *App) DisconnectWebsocket(client *msgsrvc.Client) {
	app.messagingService.Disconnect <- client

	close(client.MessagingChannel)
}

func (app *App) JoinChannels(clientObj *msgsrvc.Client, serverId uuid.UUID, channels []uuid.UUID) {
	app.messagingService.JoinChannels(clientObj, serverId, channels)
}

func (app *App) QuitChannel(clientObj *msgsrvc.Client, channelId uuid.UUID) {
	app.messagingService.QuitChannel(clientObj, channelId)
}

func (app *App) QuitServer(clientObj *msgsrvc.Client, serverId uuid.UUID) {
	channels, _ := app.GetServerChannels(serverId, clientObj.ID, 0, 1000)

	for _, channel := range *channels {
		app.QuitChannel(clientObj, channel.ID)
	}

	app.messagingService.QuitServer(clientObj, serverId)
}

func (app *App) RemoveChannel(channelId uuid.UUID) {
	app.messagingService.RemoveChannel(channelId)
}

func (app *App) RemoveServer(serverId uuid.UUID) {
	app.messagingService.RemoveServer(serverId)
}

func (app *App) GetServerMemberRole(serverId, userId uuid.UUID) (string, error) {
	return app.db.GetServerMemberRole(serverId, userId)
}
