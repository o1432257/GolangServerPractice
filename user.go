package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 創建一個用戶的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,

		server: server,
	}

	//啟動監聽當前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 用戶上線
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//廣播當前用戶上線了
	u.server.BoardCast(u, "已上線")
}

// 用戶下線
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	//廣播當前用戶上線了
	u.server.BoardCast(u, "已下線")
}

func (u *User) SendMessage(msg string) {
	u.conn.Write([]byte(msg))
}

// 處理訊息
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查詢當前在線用戶
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在線...\n"
			u.SendMessage(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else {
		u.server.BoardCast(u, msg)
	}
}

// 監聽當前user channel 一旦有消息 發送訊息到客戶端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
