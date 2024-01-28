package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/mohamed-sawy/critch-backend/internal/application/application"
	"github.com/mohamed-sawy/critch-backend/internal/application/core/msgsrvc"
	"log"
	"net/http"
)

type connection struct {
	clientObj           *msgsrvc.Client
	websocketConnection *websocket.Conn
}

func (api *Adapter) connectWebsocket(ctx *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	websocketConnection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	clientID, _ := ctx.Get("user_id")

	clientObj, err := api.app.ConnectWebsocket(clientID.(uuid.UUID))
	if err != nil {
		reportWebsocketError(websocketConnection, err)
		websocketConnection.Close()
		return
	}

	client := &connection{
		clientObj:           clientObj,
		websocketConnection: websocketConnection,
	}

	go receiveMessages(client, api.app)
	go sendMessages(client, api.app)
}

func receiveMessages(client *connection, app application.AppI) {
	defer func() {
		client.websocketConnection.Close()
	}()

	for {
		message, ok := app.ReceiveMessages(client.clientObj)
		if !ok {
			return
		}

		err := client.websocketConnection.WriteJSON(message)
		if err != nil {
			reportWebsocketError(client.websocketConnection, err)
		}
	}
}

func sendMessages(client *connection, app application.AppI) {
	defer func() {
		app.DisconnectWebsocket(client.clientObj)
		client.websocketConnection.Close()
	}()

	for {
		message := &msgsrvc.IncomingMessage{SenderID: client.clientObj.ID}
		err := client.websocketConnection.ReadJSON(message)
		if err != nil {
			reportWebsocketError(client.websocketConnection, err)
			break
		}

		err = app.SendMessages(message)
		if err != nil {
			reportWebsocketError(client.websocketConnection, err)
		}
	}
}

func reportWebsocketError(websocketConnection *websocket.Conn, err error) {
	log.Println(err)
	websocketConnection.WriteJSON(map[string]any{
		"type":    msgsrvc.ERROR,
		"message": err.Error(),
	})
}
