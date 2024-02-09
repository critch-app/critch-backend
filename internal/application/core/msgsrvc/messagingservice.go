package msgsrvc

import (
	"github.com/critch-app/critch-backend/internal/application/core/entities"
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

			//srvc.Broadcast <- &BroadcastMessage{
			//	Type: NOTIFICATION,
			//	Message: map[string]any{
			//		"type": LOGGED_IN,
			//		"data": map[string]any{
			//			"sender_id": newClient.ClientObj.ID,
			//		},
			//	},
			//}
		case client := <-srvc.Disconnect:
			for _, server := range srvc.ServerClients {
				delete(server, client.ID)
				//srvc.Broadcast <- &BroadcastMessage{
				//	Type: NOTIFICATION,
				//	Message: map[string]any{
				//		"type": LOGGED_OUT,
				//		"data": map[string]any{
				//			"sender_id": client.ID,
				//		},
				//	},
				//}
			}

			for _, channel := range srvc.ChannelClients {
				delete(channel, client.ID)
			}
		case message := <-srvc.Broadcast:
			if message.Type == NOTIFICATION {
				server := srvc.ServerClients[message.ServerId]
				for _, client := range server {
					client.MessagingChannel <- message.Message
				}
			} else if message.Type == MESSAGE {
				outgoingMessage := map[string]any{
					"server_id": message.ServerId,
				}
				if message.ServerId == uuid.Nil {
					outgoingMessage["server_id"] = nil
				}

				switch message.Message.(type) {
				case *entities.ServerMessage:
					messageModel := message.Message.(*entities.ServerMessage)
					outgoingMessage["id"] = messageModel.ID
					outgoingMessage["channel_id"] = messageModel.ChannelID
					outgoingMessage["sender_id"] = messageModel.SenderID
					outgoingMessage["content"] = messageModel.Content
					outgoingMessage["attachment"] = messageModel.Attachment
					outgoingMessage["sent_at"] = messageModel.SentAt
					outgoingMessage["updated_at"] = messageModel.UpdatedAt
				case *entities.DirectMessage:
					messageModel := message.Message.(*entities.DirectMessage)
					outgoingMessage["id"] = messageModel.ID
					outgoingMessage["channel_id"] = messageModel.ChannelID
					outgoingMessage["sender_id"] = messageModel.SenderID
					outgoingMessage["content"] = messageModel.Content
					outgoingMessage["attachment"] = messageModel.Attachment
					outgoingMessage["sent_at"] = messageModel.SentAt
					outgoingMessage["updated_at"] = messageModel.UpdatedAt
				}

				channel := srvc.ChannelClients[message.ChannelId]
				for _, client := range channel {
					client.MessagingChannel <- map[string]any{
						"type": MESSAGE,
						"data": outgoingMessage,
					}
				}
			}
		}
	}
}

func (srvc *MessagingService) JoinChannels(clientObj *Client, serverId uuid.UUID, channels []uuid.UUID) {
	if srvc.ServerClients[serverId] == nil {
		srvc.ServerClients[serverId] = make(map[uuid.UUID]*Client)
	}

	srvc.ServerClients[serverId][clientObj.ID] = clientObj

	for _, channelId := range channels {
		if srvc.ChannelClients[channelId] == nil {
			srvc.ChannelClients[channelId] = make(map[uuid.UUID]*Client)
		}

		srvc.ChannelClients[channelId][clientObj.ID] = clientObj
	}
}

func (srvc *MessagingService) QuitChannel(clientObj *Client, channelId uuid.UUID) {
	delete(srvc.ChannelClients[channelId], clientObj.ID)
}

func (srvc *MessagingService) QuitServer(clientObj *Client, serverId uuid.UUID) {
	delete(srvc.ServerClients[serverId], clientObj.ID)
}

func (srvc *MessagingService) RemoveChannel(channelId uuid.UUID) {
	delete(srvc.ChannelClients, channelId)
}

func (srvc *MessagingService) RemoveServer(serverId uuid.UUID) {
	delete(srvc.ServerClients, serverId)
}
