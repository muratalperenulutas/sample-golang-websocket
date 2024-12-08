package controller

import "github.com/gorilla/websocket"

type ClientController struct {
	clients map[*websocket.Conn]bool
}

func NewClientController() *ClientController {
	return &ClientController{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (c *ClientController) AddClient(conn *websocket.Conn, isIotDevice bool) {
	c.clients[conn] = isIotDevice
}

func (c *ClientController) RemoveClient(conn *websocket.Conn) {
	delete(c.clients, conn)
}

func (c *ClientController) GetConnection(conn *websocket.Conn) (bool, bool) {
	isIotDevice, exists := c.clients[conn]
	return isIotDevice, exists
}

func (c *ClientController) GetAllConnections() map[*websocket.Conn]bool {
	return c.clients
}
