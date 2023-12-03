package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamed-sawy/critch-backend/internal/application/application"
)

type Adapter struct {
	app    application.AppI
	router *gin.Engine
}

func NewAdapter(app application.AppI) *Adapter {
	router := gin.Default()

	api := &Adapter{app: app, router: router}
	api.setupRouting()

	return api
}

func (api *Adapter) Run() error {
	return api.router.Run()
}

func (api *Adapter) setupRouting() {
	v1 := api.router.Group("/v1")

	v1.POST("/users", api.signup)
	v1.POST("/login", api.login)
	v1.GET("/users", api.getAllUsers)
	v1.GET("/users/:user-id", api.getUser)
	v1.DELETE("/users/:user-id", api.deleteUser)
	v1.PATCH("/users/:user-id", api.updateUser)
	v1.GET("/users/:user-id/servers", api.getUserServers)
	v1.GET("/users/:user-id/channels", api.getUserChannels)

	v1.GET("/servers", api.getAllServers)
	v1.POST("/servers", api.createServer)
	v1.GET("/servers/:server-id", api.getServer)
	v1.DELETE("/servers/:server-id", api.deleteServer)
	v1.PATCH("/servers/:server-id", api.updateServer)
	v1.GET("/servers/:server-id/users", api.getServerMembers)
	v1.PUT("/servers/:server-id/users/:user-id", api.addServerMember)
	v1.DELETE("/servers/:server-id/users/:user-id", api.removeServerMember)
	v1.GET("/servers/:server-id/channels", api.getServerChannels)

	v1.GET("/channels", api.getAllChannels)
	v1.POST("/channels", api.createChannel)
	v1.GET("/channels/:channel-id", api.getChannel)
	v1.DELETE("/channels/:channel-id", api.deleteChannel)
	v1.PATCH("/channels/:channel-id", api.updateChannel)
	v1.GET("/channels/:channel-id/users", api.getChannelMembers)
	v1.PUT("/channels/:channel-id/users/:user-id", api.addChannelMember)
	v1.DELETE("/channels/:channel-id/users/:user-id", api.removeChannelMember)
	v1.GET("/channels/:channel-id/messages", api.getChannelMessages)

	v1.GET("/messages/:message-id", api.getMessage)
	v1.DELETE("/messages/:message-id", api.deleteMessage)
	v1.PATCH("/messages/:message-id", api.updateMessage)
}
