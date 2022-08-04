package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 創建一個server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

func (s *Server) Handler(conn net.Conn) {
	//當前連接的業務
	fmt.Println("連結建立成功")
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
