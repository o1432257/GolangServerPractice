package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	user := NewUser(conn, s)
	//用戶上線了
	user.Online()

	//監聽用戶是否活躍的頻道
	isLive := make(chan bool)

	//接受客戶端發送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用戶的消息
			msg := string(buf[:n-1])

			//將用戶的消息進行廣播
			user.DoMessage(msg)

			//用戶的任意消息,代表是活躍的
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
			//當前用戶是活躍的 重製定時器
			//不做任何事 只激活select

		case <-time.After(time.Minute * 10):
			//已經超時,將當前的User強制關閉
			user.SendMessage("已被踢出")

			//關閉通道
			close(user.C)

			//關閉連結
			conn.Close()

			//退出當前Handler
			return
		}
	}
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
