package api

import (
	"errors"
	"github.com/gin-gonic/gin"
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

	token, exists := ctx.GetQuery("token")
	if !exists {
		reportError(ctx, http.StatusBadRequest, errors.New("unauthorized"))
		return
	}

	websocketConnection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		reportError(ctx, http.StatusInternalServerError, err)
		return
	}

	clientId, err := api.app.ValidateJWTToken(token)
	if err != nil {
		reportWebsocketError(websocketConnection, err)
		websocketConnection.Close()
		return
	}

	clientObj, err := api.app.ConnectWebsocket(clientId)
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
