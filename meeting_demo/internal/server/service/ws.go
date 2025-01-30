package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type WsP2PConnRequest struct {
	RoomIdentity string `json:"room_id" uri:"room_id"`
	UserIdentity string `json:"user_id" uri:"user_id"`
}

type WsP2PMessage struct {
	RoomIdentity string `json:"room_id"`
	UserIdentity string `json:"user_id"`
	Key          string `json:"key"` // 表明传递信息为sdp还是candidate
	Value        any    `json:"value"`
}

// 通过room_id和user_id找到对应的连接
// 存储房间对应的用户信息，键为room_id，值为用户信息map（user_id为键、user connection为值）
var roomUserMap = sync.Map{}

func WsP2PConnection(c *gin.Context) {
	wsReq := new(WsP2PConnRequest)
	err := c.ShouldBindUri(wsReq)
	if err != nil {
		log.Println("bind ws request failed: ", err)
		return
	}

	// 将HTTP连接升级为WebSocket连接
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	conn, err1 := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err1 != nil {
		log.Println("create websocket failed: ", err1)
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	// 建立user_id和conn的map，并将该map存入roomUserMap中
	userConnMap, _ := roomUserMap.LoadOrStore(wsReq.RoomIdentity, &sync.Map{})
	userConn := userConnMap.(*sync.Map) // 将返回的any类型断言为*sync.Map
	userConn.Store(wsReq.UserIdentity, conn)

	for {
		// 监听消息，另一端会将sdp或者candidate信息通过websocket发送过来
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read message from websocket failed: ", err)
			return
		}

		message := new(WsP2PMessage)
		err = json.Unmarshal(data, message)
		if err != nil {
			log.Println("unmarshal websocket message failed: ", err)
			continue
		}

		// 发送消息到一个房间内的所有用户
		users, ok := roomUserMap.Load(wsReq.RoomIdentity)
		if ok {
			users.(*sync.Map).Range(func(key, value any) bool {
				//if key == message.UserIdentity {
				//    return true
				//}
				err := value.(*websocket.Conn).WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println("write message to websocket failed: ", err)
					return false
				}
				return true
			})
		}
	}
}
