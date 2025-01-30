package main

import (
	"bufio"
	"fmt"
	"github.com/pion/webrtc/v4"
	"log"
	"meeting_demo/utils"
	"os"
	"strconv"
	"time"
)

func main() {
	// 1. create peer connection
	peerConn, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Println("create peer peerConn error: ", err)
		return
	}

	// 2. listen remote data channel
	// 当远端创建一个新的Data Channel并添加到现有的PeerConnection中时触发
	peerConn.OnDataChannel(func(channel *webrtc.DataChannel) {
		// 当Data Channel成功打开并准备好收发时触发
		channel.OnOpen(func() {
			log.Println("answer data channel open")
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()

			cnt := 0
			for range ticker.C {
				err := channel.SendText("hello yxj " + strconv.Itoa(cnt))
				if err != nil {
					log.Println("send data failed: ", err)
					return
				}
				cnt++
			}
		})

		// 当新消息从远端到达Data Channel时触发
		channel.OnMessage(func(msg webrtc.DataChannelMessage) {
			if msg.IsString {
				fmt.Println("receive string message: ", string(msg.Data))
			} else {
				fmt.Println("receive binary message with length: ", len(msg.Data))
			}
		})
	})

	// 3. waiting for inputting offer string
	fmt.Println("Please input another computer's offer:")
	readString, err1 := bufio.NewReader(os.Stdin).ReadString('\n')
	if err1 != nil {
		log.Println("input offer failed: ", err1)
		return
	}
	var offer webrtc.SessionDescription
	utils.Decode(readString, &offer)

	// 4. set remote description
	err = peerConn.SetRemoteDescription(offer)
	if err != nil {
		log.Println("set remote description failed: ", err)
		return
	}

	// 5. create answer
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		log.Println("create answer failed: ", err)
		return
	}

	// 6. set local description
	err = peerConn.SetLocalDescription(answer)
	if err != nil {
		log.Println("set local description failed: ", err)
		return
	}

	// 7. gather complete
	// 等待ICE候选者收集完成
	//（在尝试建立 WebRTC 连接时，PeerConnection 需要找到所有可能的通信路径，这通常涉及到使用 STUN 或 TURN 服务器来穿透 NAT 和防火墙，这个过程被称为 ICE 候选者收集）
	gatherComplete := webrtc.GatheringCompletePromise(peerConn)
	<-gatherComplete

	// 8. print answer(change sdp)
	fmt.Println("Answer: ", utils.Encode(peerConn.LocalDescription()))

	select {}
}
