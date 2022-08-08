package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {

	//創建客戶端對象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	//連結server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn

	//返回對象
	return client
}

// 處理server回應的消息 直接顯示道標準輸出
func (c *Client) DealResponse() {
	//一旦 client.conn 有消息,就直接copy到stdout標準輸出上, 阻塞監聽
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.公開聊天室")
	fmt.Println("2.私人聊天室")
	fmt.Println("3.更新用戶名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println(">>>>>請輸入合法範圍內的數字<<<<<")
		return false
	}

}

func (c *Client) UpdateName() bool {

	fmt.Println(">>>>>請輸入用戶名:")
	fmt.Scanln(&c.Name)

	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {
		}

		//根據不同的模式處裡不同的業務
		switch c.flag {
		case 1:
			//公開聊天室
			c.PublicChat()
		case 2:
			//私人聊天室
			c.PrivateChat()
		case 3:
			//更新用戶名
			c.UpdateName()
			break
		}
	}
}

func (c *Client) PublicChat() {
	//提示用戶輸入消息
	var chatMsg string

	fmt.Println(">>>>>請輸入聊天內容,exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := c.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>>請輸入聊天內容,exit退出")
		fmt.Scanln(&chatMsg)
	}
}

// 查詢在線用戶
func (c *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write, err")
		return
	}
}

// 私人聊天室
func (c *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	c.SelectUsers()
	fmt.Println(">>>>>請輸入聊天對象的用戶名, exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>請輸入消息內容, exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := c.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>>請輸入消息內容, exit退出")
			fmt.Scanln(&chatMsg)
		}

		c.SelectUsers()
		fmt.Println(">>>>>請輸入聊天對象的用戶名, exit退出")
		fmt.Scanln(&remoteName)
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "設置服務器IP地址(默認127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "設置服務器端口(默認8888)")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>連接服務器失敗")
		return
	}

	//單獨開啟一個goroutine, 去處理server返回的消息
	go client.DealResponse()

	fmt.Println(">>>>>>連接服務器成功")

	//啟動客戶端的業務
	client.Run()
}
