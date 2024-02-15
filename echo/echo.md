# echo

## 依赖包
 - ```github.com/gorilla/websocket```
## websocket
- webSocket协议
 - 传入的http协议会转换成webSocket协议
  - ```var upgrader = websocket.Upgrader{} // use default options```
  - ```conn, err := upgrader.Upgrade(w, r, nil)```
- 虽然webSocket的客户端（一般不会用go写，但是贴出代码方便理解)：<br>
```
func main() {
	// 连接WebSocket服务器
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}

	// 读取消息
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Received message:", string(p), "MessageType:", messageType)
}
```
- websocket不允许用户传入header参数。如果你想传递额外信息，官方建议就是三种：URL、cookie