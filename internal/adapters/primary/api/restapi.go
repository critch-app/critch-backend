package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mohamed-sawy/critch-backend/internal/application/application"
)

type Adapter struct {
	app    application.AppI
	router *gin.Engine
}

func NewAdapter(app application.AppI) *Adapter {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	router.Use(cors.New(config))

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

	authorized := v1.Group("/", api.authenticate)

	authorized.GET("/users", api.getAllUsers)
	authorized.GET("/users/:user-id", api.getUser)
	authorized.DELETE("/users/:user-id", api.deleteUser)
	authorized.PATCH("/users/:user-id", api.updateUser)
	authorized.GET("/users/:user-id/servers", api.getUserServers)
	authorized.GET("/users/:user-id/channels", api.getUserChannels)

	authorized.GET("/servers", api.getAllServers)
	authorized.POST("/servers", api.createServer)
	authorized.GET("/servers/:server-id", api.getServer)
	authorized.DELETE("/servers/:server-id", api.deleteServer)
	authorized.PATCH("/servers/:server-id", api.updateServer)
	authorized.GET("/servers/:server-id/users", api.getServerMembers)
	authorized.PUT("/servers/:server-id/users/:user-id", api.addServerMember)
	authorized.DELETE("/servers/:server-id/users/:user-id", api.removeServerMember)
	authorized.GET("/servers/:server-id/channels", api.getServerChannels)

	authorized.GET("/channels", api.getAllChannels)
	authorized.POST("/channels", api.createChannel)
	authorized.GET("/channels/:channel-id", api.getChannel)
	authorized.DELETE("/channels/:channel-id", api.deleteChannel)
	authorized.PATCH("/channels/:channel-id", api.updateChannel)
	authorized.GET("/channels/:channel-id/users", api.getChannelMembers)
	authorized.PUT("/channels/:channel-id/users/:user-id", api.addChannelMember)
	authorized.DELETE("/channels/:channel-id/users/:user-id", api.removeChannelMember)
	authorized.GET("/channels/:channel-id/messages", api.getChannelMessages)

	authorized.GET("/messages/:message-id", api.getMessage)
	authorized.DELETE("/messages/:message-id", api.deleteMessage)
	authorized.PATCH("/messages/:message-id", api.updateMessage)
}
