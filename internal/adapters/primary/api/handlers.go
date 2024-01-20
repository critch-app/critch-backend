package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/entities"
	"log"
	"net/http"
	"strconv"
)

func (api *Adapter) login(ctx *gin.Context) {
	credentials := &loginRequest{}

	err := ctx.ShouldBindJSON(credentials)
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	token, userId, err := api.app.Login(credentials.Email, credentials.Password)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"token":   token,
		"user_id": userId,
	})
}

func (api *Adapter) signup(ctx *gin.Context) {
	user := &entities.User{}

	err := ctx.ShouldBindJSON(user)
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.Signup(user)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, getResponseUser(user))
}

func (api *Adapter) getAllUsers(ctx *gin.Context) {
	offset, limit := getPagination(ctx)

	users, err := api.app.GetAllUsers(offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	usersData := make([]gin.H, len(*users))
	for idx, user := range *users {
		usersData[idx] = getResponseUser(&user)
	}

	ctx.JSON(http.StatusOK, usersData)
}

func (api *Adapter) getUser(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := api.app.GetUser(userId)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseUser(user))
}

func (api *Adapter) deleteUser(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.DeleteUser(userId)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) updateUser(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	user := &entities.User{}

	err = ctx.ShouldBindJSON(user)
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	user.ID = userId

	err = api.app.UpdateUser(user)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseUser(user))
}

func (api *Adapter) getUserServers(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	offset, limit := getPagination(ctx)

	servers, err := api.app.GetUserServers(userId, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	serversData := make([]gin.H, len(*servers))
	for idx, server := range *servers {
		serversData[idx] = gin.H{
			"server_id": server.ServerID,
			"joined_at": server.JoinedAt,
			"role":      server.Role,
		}
	}

	ctx.JSON(http.StatusOK, serversData)
}

func (api *Adapter) getUserChannels(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	offset, limit := getPagination(ctx)

	channels, err := api.app.GetUserDMChannels(userId, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	channelsData := make([]gin.H, len(*channels))
	for idx, channel := range *channels {
		channelsData[idx] = gin.H{
			"channel_id": channel.ChannelID,
		}
	}

	ctx.JSON(http.StatusOK, channelsData)
}

func (api *Adapter) getAllServers(ctx *gin.Context) {
	offset, limit := getPagination(ctx)

	servers, err := api.app.GetAllServers(offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	serversData := make([]gin.H, len(*servers))
	for idx, server := range *servers {
		serversData[idx] = getResponseServer(&server)
	}

	ctx.JSON(http.StatusOK, serversData)
}

func (api *Adapter) createServer(ctx *gin.Context) {
	serverRequest := &createServerRequest{}

	err := ctx.ShouldBindJSON(serverRequest)
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	OwnerID, err := uuid.Parse(serverRequest.OwnerID)
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	server := &entities.Server{
		Name:        serverRequest.Name,
		Description: serverRequest.Description,
		Photo:       serverRequest.Photo,
	}

	err = api.app.CreateServer(server, OwnerID)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, getResponseServer(server))
}

func (api *Adapter) getServer(ctx *gin.Context) {
	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	server, err := api.app.GetServer(serverId)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseServer(server))
}

func (api *Adapter) deleteServer(ctx *gin.Context) {
	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.DeleteServer(serverId)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) updateServer(ctx *gin.Context) {
	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	server := &entities.Server{}

	err = ctx.ShouldBindJSON(server)
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	server.ID = serverId

	err = api.app.UpdateServer(server)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseServer(server))
}

func (api *Adapter) getServerMembers(ctx *gin.Context) {
	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	offset, limit := getPagination(ctx)

	members, err := api.app.GetServerMembers(serverId, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	membersData := make([]gin.H, len(*members))
	for idx, member := range *members {
		membersData[idx] = gin.H{
			"user_id":   member.UserID,
			"joined_at": member.JoinedAt,
			"role":      member.Role,
		}
	}

	ctx.JSON(http.StatusOK, membersData)
}

func (api *Adapter) addServerMember(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.AddServerMember(serverId, userId)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) removeServerMember(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.RemoveServerMember(serverId, userId)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) getServerChannels(ctx *gin.Context) {
	serverId, err := uuid.Parse(ctx.Param("server-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	offset, limit := getPagination(ctx)

	channels, err := api.app.GetServerChannels(serverId, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	channelsData := make([]gin.H, len(*channels))
	for idx, channel := range *channels {
		channelsData[idx] = gin.H{
			"channel_id": channel.ID,
		}
	}

	ctx.JSON(http.StatusOK, channelsData)
}

func (api *Adapter) getAllChannels(ctx *gin.Context) {
	offset, limit := getPagination(ctx)

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channels any
	if isServerChannel {
		channels = &[]entities.ServerChannel{}
	} else {
		channels = &[]entities.DMChannel{}
	}

	err := api.app.GetAllChannels(channels, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	channelsData := getResponseChannelArray(channels, isServerChannel)

	ctx.JSON(http.StatusOK, channelsData)
}

func (api *Adapter) createChannel(ctx *gin.Context) {
	_, isServerChannel := ctx.GetQuery("isServerChannel")

	var (
		err     error
		channel any
	)

	if isServerChannel {
		channel = &entities.ServerChannel{}
		err = ctx.ShouldBindJSON(channel)
	} else {
		channel = &entities.DMChannel{}
		err = ctx.ShouldBindJSON(channel)
	}

	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.CreateChannel(channel)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, getResponseChannel(channel, isServerChannel))
}

func (api *Adapter) getChannel(ctx *gin.Context) {
	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channel any
	if isServerChannel {
		channel = &entities.ServerChannel{Channel: entities.Channel{ID: channelId}}
	} else {
		channel = &entities.DMChannel{Channel: entities.Channel{ID: channelId}}
	}

	err = api.app.GetChannel(channel)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseChannel(channel, isServerChannel))
}

func (api *Adapter) deleteChannel(ctx *gin.Context) {
	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channel any
	if isServerChannel {
		channel = &entities.ServerChannel{Channel: entities.Channel{ID: channelId}}
	} else {
		channel = &entities.DMChannel{Channel: entities.Channel{ID: channelId}}
	}

	err = api.app.DeleteChannel(channel)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) updateChannel(ctx *gin.Context) {
	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channel any
	if isServerChannel {
		channelModel := &entities.ServerChannel{}
		err = ctx.ShouldBindJSON(channelModel)
		channelModel.ID = channelId
		channel = channelModel
	} else {
		channelModel := &entities.DMChannel{}
		err = ctx.ShouldBindJSON(channelModel)
		channelModel.ID = channelId
		channel = channelModel
	}

	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.UpdateChannel(channel)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseChannel(channel, isServerChannel))
}

func (api *Adapter) getChannelMembers(ctx *gin.Context) {
	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	offset, limit := getPagination(ctx)

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channelMembers any
	if isServerChannel {
		channelMembers = &[]entities.ServerChannelMember{}
	} else {
		channelMembers = &[]entities.DMChannelMember{}
	}

	err = api.app.GetChannelMembers(channelMembers, channelId, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	membersData := getResponseChannelMembersArray(channelMembers, isServerChannel)

	ctx.JSON(http.StatusOK, membersData)
}

func (api *Adapter) addChannelMember(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channelMember any
	if isServerChannel {
		channelMember = &entities.ServerChannelMember{
			ChannelID: channelId,
			UserID:    userId,
		}
	} else {
		channelMember = &entities.DMChannelMember{
			ChannelID: channelId,
			UserID:    userId,
		}
	}

	err = api.app.AddChannelMember(channelMember)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) removeChannelMember(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Param("user-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channelMember any
	if isServerChannel {
		channelMember = &entities.ServerChannelMember{
			ChannelID: channelId,
			UserID:    userId,
		}
	} else {
		channelMember = &entities.DMChannelMember{
			ChannelID: channelId,
			UserID:    userId,
		}
	}

	err = api.app.RemoveChannelMember(channelMember)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) getChannelMessages(ctx *gin.Context) {
	channelId, err := uuid.Parse(ctx.Param("channel-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	offset, limit := getPagination(ctx)

	_, isServerChannel := ctx.GetQuery("isServerChannel")
	var channelMessages any
	if isServerChannel {
		channelMessages = &[]entities.ServerMessage{}
	} else {
		channelMessages = &[]entities.DirectMessage{}
	}

	err = api.app.GetChannelMessages(channelMessages, channelId, offset, limit)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	membersData := getResponseMessageArray(channelMessages, isServerChannel)

	ctx.JSON(http.StatusOK, membersData)
}

func (api *Adapter) getMessage(ctx *gin.Context) {
	messageId, err := uuid.Parse(ctx.Param("message-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerMessage := ctx.GetQuery("isServerMessage")
	var message any
	if isServerMessage {
		message = &entities.ServerMessage{Message: entities.Message{ID: messageId}}
	} else {
		message = &entities.DirectMessage{Message: entities.Message{ID: messageId}}
	}

	err = api.app.GetMessage(message)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseMessage(message, isServerMessage))
}

func (api *Adapter) deleteMessage(ctx *gin.Context) {
	messageId, err := uuid.Parse(ctx.Param("message-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerMessage := ctx.GetQuery("isServerMessage")
	var message any
	if isServerMessage {
		message = &entities.ServerMessage{Message: entities.Message{ID: messageId}}
	} else {
		message = &entities.DirectMessage{Message: entities.Message{ID: messageId}}
	}

	err = api.app.DeleteMessage(message)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (api *Adapter) updateMessage(ctx *gin.Context) {
	messageId, err := uuid.Parse(ctx.Param("message-id"))
	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	_, isServerMessage := ctx.GetQuery("isServerMessage")
	var message any
	if isServerMessage {
		messageModel := &entities.ServerMessage{}
		err = ctx.ShouldBindJSON(messageModel)
		messageModel.ID = messageId
		message = messageModel
	} else {
		messageModel := &entities.DirectMessage{}
		err = ctx.ShouldBindJSON(messageModel)
		messageModel.ID = messageId
		message = messageModel
	}

	if err != nil {
		reportError(ctx, http.StatusBadRequest, err)
		return
	}

	err = api.app.UpdateMessage(message)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getResponseMessage(message, isServerMessage))
}

func getPagination(ctx *gin.Context) (offset int, limit int) {
	const MIN_OFFSET = 0
	const DEFAULT_OFFSET = 0
	const MIN_LIMIT = 1
	const MAX_LIMIT = 100
	const DEFAULT_LIMIT = 50

	offset = DEFAULT_OFFSET
	limit = DEFAULT_LIMIT

	offsetQuery := ctx.Query("offset")
	limitQuery := ctx.Query("limit")

	offsetValue, err := strconv.Atoi(offsetQuery)
	if err == nil {
		offset = offsetValue
	}
	limitValue, err := strconv.Atoi(limitQuery)
	if err == nil {
		limit = limitValue
	}

	if offset < MIN_OFFSET {
		offset = DEFAULT_OFFSET
	}
	if limit < MIN_LIMIT || limit > MAX_LIMIT {
		limit = DEFAULT_LIMIT
	}

	return offset, limit
}

func getResponseUser(user *entities.User) gin.H {
	return gin.H{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"phone":      user.Phone,
		"status":     user.Status,
		"photo":      user.Photo,
		"time_zone":  user.TimeZone,
		"created_at": user.CreatedAt,
		"last_seen":  user.LastSeen,
	}
}

func getResponseServer(server *entities.Server) gin.H {
	return gin.H{
		"id":          server.ID,
		"name":        server.Name,
		"description": server.Description,
		"photo":       server.Photo,
		"created_at":  server.CreatedAt,
	}
}

func getResponseChannel(channel any, isServerChannel bool) gin.H {
	if isServerChannel {
		channelModel := channel.(*entities.ServerChannel)
		return gin.H{
			"id":          channelModel.ID,
			"server_id":   channelModel.CreatedAt,
			"name":        channelModel.Name,
			"description": channelModel.Description,
			"created_at":  channelModel.CreatedAt,
		}
	}

	channelModel := channel.(*entities.DMChannel)
	return gin.H{
		"id":          channelModel.ID,
		"name":        channelModel.Name,
		"description": channelModel.Description,
		"created_at":  channelModel.CreatedAt,
	}
}

func getResponseChannelArray(channels any, isServerChannel bool) []gin.H {
	if isServerChannel {
		channelsArray := channels.(*[]entities.ServerChannel)
		channelsData := make([]gin.H, len(*channelsArray))
		for idx, channel := range *channelsArray {
			channelsData[idx] = getResponseChannel(&channel, isServerChannel)
		}

		return channelsData
	}

	channelsArray := channels.(*[]entities.DMChannel)
	channelsData := make([]gin.H, len(*channelsArray))
	for idx, channel := range *channelsArray {
		channelsData[idx] = getResponseChannel(&channel, isServerChannel)
	}

	return channelsData
}

func getResponseChannelMembersArray(channelMembers any, isServerChannel bool) []gin.H {
	if isServerChannel {
		membersArray := channelMembers.(*[]entities.ServerChannelMember)
		membersData := make([]gin.H, len(*membersArray))
		for idx, member := range *membersArray {
			membersData[idx] = gin.H{
				"user_id": member.UserID,
			}
		}

		return membersData
	}

	membersArray := channelMembers.(*[]entities.DMChannelMember)
	membersData := make([]gin.H, len(*membersArray))
	for idx, member := range *membersArray {
		membersData[idx] = gin.H{
			"user_id": member.UserID,
		}
	}

	return membersData
}

func getResponseMessage(message any, isServerMessage bool) gin.H {
	if isServerMessage {
		messageModel := message.(*entities.ServerMessage)
		return gin.H{
			"id":         messageModel.ID,
			"content":    messageModel.Content,
			"attachment": messageModel.Attachment,
			"sent_at":    messageModel.SentAt,
			"updated_at": messageModel.UpdatedAt,
		}
	}

	messageModel := message.(*entities.DirectMessage)
	return gin.H{
		"id":         messageModel.ID,
		"content":    messageModel.Content,
		"attachment": messageModel.Attachment,
		"sent_at":    messageModel.SentAt,
		"updated_at": messageModel.UpdatedAt,
	}
}

func getResponseMessageArray(channelMessages any, isServerMessage bool) []gin.H {
	if isServerMessage {
		messagesArray := channelMessages.(*[]entities.ServerMessage)
		messagesData := make([]gin.H, len(*messagesArray))
		for idx, message := range *messagesArray {
			messagesData[idx] = getResponseMessage(message, isServerMessage)
		}

		return messagesData
	}

	messagesArray := channelMessages.(*[]entities.DirectMessage)
	messagesData := make([]gin.H, len(*messagesArray))
	for idx, message := range *messagesArray {
		messagesData[idx] = getResponseMessage(message, isServerMessage)
	}

	return messagesData
}

func reportError(ctx *gin.Context, errorCode int, err error) {
	log.Println(err)
	ctx.JSON(errorCode, gin.H{
		"message": err.Error(),
	})
}
