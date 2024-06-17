package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type User struct {
	uid string
	c   *websocket.Conn
}

type Server struct {
	conns map[string]*User
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var server = Server{conns: make(map[string]*User)}

func main() {
	// 创建Gin应用
	app := gin.Default()

	// 注册WebSocket路由
	app.GET("/ws", WebSocketHandler)

	// 启动应用
	err := app.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func WebSocketHandler(c *gin.Context) {
	uid := c.GetHeader("id")
	if user, ok := server.conns[uid]; ok {
		user.c.Close()
	}
	// 获取WebSocket连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	server.conns[uid] = &User{
		uid: uid,
		c:   conn,
	}

	// 处理WebSocket消息
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("messageType:", messageType)
		fmt.Println("data:", string(data))

		// 输出WebSocket消息内容
		conn.WriteMessage(websocket.TextMessage, data)
	}

	// 关闭WebSocket连接
	conn.Close()
}
