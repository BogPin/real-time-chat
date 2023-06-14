package wss

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)

type room struct {
	mu    sync.RWMutex
	Id    int
	conns map[int]*websocket.Conn
}

func (r *room) Add(userId int, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[userId] = conn
}

func (r *room) Remove(userId int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, userId)
}

func (r *room) Send(fromUser int, msg Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for userId, conn := range r.conns {
		if userId != fromUser {
			if err := conn.WriteJSON(msg); err != nil {
				fmt.Printf("error while sending msg %v in room %v from user %v:\n", msg, r.Id, fromUser)
			}
		}
	}
}

type safeRooms struct {
	mu    sync.RWMutex
	rooms []*room
}

func (sr *safeRooms) Get(roomId int) (*room, error) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	i := slices.IndexFunc(sr.rooms, func(r *room) bool { return r.Id == roomId })
	if i != -1 {
		return sr.rooms[i], nil
	}
	return nil, fmt.Errorf("no room with id %d", roomId)
}

func (sr *safeRooms) Add(roomId int) *room {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	room := &room{
		Id:    roomId,
		conns: make(map[int]*websocket.Conn),
	}
	sr.rooms = append(sr.rooms, room)
	return room
}

func (sr *safeRooms) Remove(roomId int) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	i := slices.IndexFunc(sr.rooms, func(r *room) bool { return r.Id == roomId })
	if i != -1 {
		sr.rooms = slices.Delete(sr.rooms, i, i+1)
		return nil
	}
	return fmt.Errorf("no room with id %d", roomId)
}

func (sr *safeRooms) GetAllForUser(userId int) []*room {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	rooms := make([]*room, 0)
	for _, room := range sr.rooms {
		if _, ok := room.conns[userId]; ok {
			rooms = append(rooms, room)
		}
	}
	return rooms
}
