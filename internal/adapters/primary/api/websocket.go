package api

import (
	"encoding/json"
	"errors"
	"github.com/critch-app/critch-backend/internal/application/application"
	"github.com/critch-app/critch-backend/internal/application/core/msgsrvc"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
		wsMessage := &webSocketMessage{}
		err := client.websocketConnection.ReadJSON(wsMessage)
		if err != nil {
			reportWebsocketError(client.websocketConnection, err)
			return
		}

		switch wsMessage.MessageType {
		case msgsrvc.MESSAGE:
			message := &msgsrvc.IncomingMessage{}

			err = json.Unmarshal(wsMessage.Data, &message)

			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}

			message.SenderId = client.clientObj.ID
			err = app.SendMessages(message)
			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}
		case msgsrvc.JOIN_CHANNEL:
			message := &msgsrvc.JoinChannel{}

			err = json.Unmarshal(wsMessage.Data, &message)

			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}

			message.SenderId = client.clientObj.ID

			app.JoinChannels(client.clientObj, message.ServerId, message.Channels)

			err = app.SendNotification(wsMessage)
			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}
		case msgsrvc.QUIT_CHANNEL:
			message := &msgsrvc.QuitChannel{}

			err = json.Unmarshal(wsMessage.Data, &message)

			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}

			message.SenderId = client.clientObj.ID

			app.QuitChannel(client.clientObj, message.ChannelId)

			err = app.SendNotification(wsMessage)
			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}
		case msgsrvc.QUIT_SERVER:
			message := &msgsrvc.QuitServer{}

			err = json.Unmarshal(wsMessage.Data, &message)

			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}

			message.SenderId = client.clientObj.ID

			app.QuitServer(client.clientObj, message.ServerId)

			err = app.SendNotification(wsMessage)
			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}
		case msgsrvc.REMOVE_CHANNEL:
			message := &msgsrvc.RemoveChannel{}

			err = json.Unmarshal(wsMessage.Data, &message)

			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}

			message.SenderId = client.clientObj.ID

			app.RemoveChannel(message.ChannelId)

			err = app.SendNotification(wsMessage)
			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}
		case msgsrvc.REMOVE_SERVER:
			message := &msgsrvc.RemoveServer{}

			err = json.Unmarshal(wsMessage.Data, &message)

			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}

			message.SenderId = client.clientObj.ID

			app.RemoveServer(message.ServerId)

			err = app.SendNotification(wsMessage)
			if err != nil {
				reportWebsocketError(client.websocketConnection, err)
				continue
			}
		default:
			log.Println("Unhandled websocket message: ", wsMessage.MessageType)
		}
	}
}

func reportWebsocketError(websocketConnection *websocket.Conn, err error) {
	log.Println(err)
	websocketConnection.WriteJSON(map[string]any{
		"type":    "error",
		"message": err.Error(),
	})
}

type webSocketMessage struct {
	MessageType string          `json:"type"  binding:"required"`
	Data        json.RawMessage `json:"data"  binding:"required"`
}
