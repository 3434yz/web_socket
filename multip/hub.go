// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	roomID string
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	//register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub(roomID string) *Hub {
	return &Hub{
		roomID:    roomID,
		broadcast: make(chan []byte),
		//register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	defer func() {
		//close(h.register)
		close(h.broadcast)
		close(h.unregister)
	}()
	for {
		select {
		//case client := <-h.register:
		//	h.clients[client] = true
		case client := <-h.unregister:
			//roomMutexes[h.roomID].Lock()
			roomMutex := roomMutexes[h.roomID]
			roomMutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			if len(h.clients) == 0 {
				fmt.Println("Delete Room", h.roomID)
				//delete(house, h.roomID)
				house.Delete(h.roomID)
				roomMutexes[h.roomID].Unlock()
				return
			}
			roomMutexes[h.roomID].Unlock()
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
