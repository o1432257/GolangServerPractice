package main

import (
	"net"
	"strings"
)

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

	//查詢當前在線用戶
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在線...\n"
			u.SendMessage(onlineMsg)
		}
		u.server.mapLock.Unlock()
		return
	}

	//更改用戶名
	if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式:rename|userName
		newName := msg[7:]

		_, ok := u.server.OnlineMap[newName]

		if ok {
			u.SendMessage("用戶名已使用")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()

			u.Name = newName
			u.SendMessage("您已更新用戶名:" + u.Name + "\n")
		}
		return
	}

	if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式:to|userName

		//取得用戶名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMessage("消息格式不正確,請使用\"to|張三|你好啊\"格式. \n")
			return
		}

		//根據用戶名，得到對方USER對象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMessage("該用戶不存在")
			return
		}

		//得到消息內容，根據對方user對象發送
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMessage("無消息內容,請重新發送\n")
			return
		}
		remoteUser.SendMessage(u.Name + "對您說" + content + "\n")
		return
	}

	//廣播
	u.server.BoardCast(u, msg)
}

// 監聽當前user channel 一旦有消息 發送訊息到客戶端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
