package msgsrvc

import (
	"github.com/google/uuid"
)

type MessagingService struct {
	ServerClients  map[uuid.UUID]map[uuid.UUID]*Client
	ChannelClients map[uuid.UUID]map[uuid.UUID]*Client
	Connect        chan *NewClient
	Disconnect     chan *Client
	Broadcast      chan *BroadcastMessage
}

func NewService() *MessagingService {
	return &MessagingService{
		ServerClients:  make(map[uuid.UUID]map[uuid.UUID]*Client),
		ChannelClients: make(map[uuid.UUID]map[uuid.UUID]*Client),
		Connect:        make(chan *NewClient),
		Disconnect:     make(chan *Client),
		Broadcast:      make(chan *BroadcastMessage, 10),
	}
}

func (srvc *MessagingService) Run() {
	for {
		select {
		case newClient := <-srvc.Connect:
			for _, serverId := range *newClient.Servers {
				if srvc.ServerClients[serverId] == nil {
					srvc.ServerClients[serverId] = make(map[uuid.UUID]*Client)
				}

				srvc.ServerClients[serverId][newClient.ClientObj.ID] = newClient.ClientObj
			}

			for _, channelId := range *newClient.Channels {
				if srvc.ChannelClients[channelId] == nil {
					srvc.ChannelClients[channelId] = make(map[uuid.UUID]*Client)
				}

				srvc.ChannelClients[channelId][newClient.ClientObj.ID] = newClient.ClientObj
			}

			srvc.Broadcast <- &BroadcastMessage{
				IsNotification: true,
				Message: map[string]any{
					"is_notification": true,
					"sender_id":       newClient.ClientObj.ID,
					"logged_in":       true,
				},
			}
		case client := <-srvc.Disconnect:
			for _, server := range srvc.ServerClients {
				delete(server, client.ID)
				srvc.Broadcast <- &BroadcastMessage{
					IsNotification: true,
					Message: map[string]any{
						"is_notification": true,
						"sender_id":       client.ID,
						"logged_out":      true,
					},
				}
			}

			for _, channel := range srvc.ChannelClients {
				delete(channel, client.ID)
			}
		case message := <-srvc.Broadcast:
			if message.IsNotification {
				for _, server := range srvc.ServerClients {
					for _, client := range server {
						client.MessagingChannel <- message.Message
					}
				}
			} else {
				channel := srvc.ChannelClients[message.ChannelId]
				for _, client := range channel {
					client.MessagingChannel <- message.Message
				}
			}
		}
	}
}
