package worker

import "github.com/gorilla/websocket"

type ConnectionTable struct {
	connMap map[string]*websocket.Conn
}

func NewConnectionTable() ConnectionTable {
	return ConnectionTable{connMap: make(map[string]*websocket.Conn)}
}

func (t ConnectionTable) Get(id string) (*websocket.Conn, bool) {
	c, ok := t.connMap[id]
	return c, ok
}

func (t ConnectionTable) Set(id string, conn *websocket.Conn) {
	t.connMap[id] = conn
}

func (t ConnectionTable) Remove(id string) {
	delete(t.connMap, id)
}
