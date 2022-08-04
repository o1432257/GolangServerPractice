package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在線用戶的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息廣播的 channel
	Message chan string
}

// 創建一個server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 監聽Message廣播消息 channel 的 goroutine , 一旦有消息發送給全部的在線 User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message

		//將msg發給全部的在線用戶
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// 廣播消息的方法
func (s *Server) BoardCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	//當前連接的業務
	fmt.Println("連結建立成功")

	user := NewUser(conn)
	//用戶上線了
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//廣播當前用戶上線了
	s.BoardCast(user, "已上線")
}

// 啟動服務器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//啟動監聽Message 的 goroutine
	go s.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go s.Handler(conn)
	}
}
