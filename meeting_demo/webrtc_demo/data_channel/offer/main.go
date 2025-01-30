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
		log.Println("create peer connection failed: ", err)
		return
	}
	defer func() {
		err := peerConn.Close()
		if err != nil {
			log.Println("close peer connection failed: ", err)
			return
		}
	}()

	// 2. create data channel
	channel, err := peerConn.CreateDataChannel("data", nil)
	if err != nil {
		return
	}

	// 当Data Channel成功打开并准备好收发时触发
	channel.OnOpen(func() {
		log.Println("offer data channel open")
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()

		cnt := 0
		for range ticker.C {
			err := channel.SendText("hello yzy " + strconv.Itoa(cnt))
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

	// 3. create offer
	offer, err1 := peerConn.CreateOffer(nil)
	if err1 != nil {
		log.Println("create offer failed: ", err1)
		return
	}

	// 4. set local description
	err = peerConn.SetLocalDescription(offer)
	if err != nil {
		log.Println("offer set local description failed: ", err)
		return
	}

	// 5. print offer(change sdp)
	fmt.Println("Offer: ", utils.Encode(offer)) // 将webrtc.SessionDescription对象转为Base64编码的JSON字符串，方便输入查看

	// 6. waiting fot inputting answer string
	fmt.Println("Please input another computer's answer: ")
	readString, err2 := bufio.NewReader(os.Stdin).ReadString('\n')
	if err2 != nil {
		log.Println("input answer description failed: ", err2)
		return
	}

	// 7. set remote description
	var answer webrtc.SessionDescription
	utils.Decode(readString, &answer)
	err = peerConn.SetRemoteDescription(answer)
	if err != nil {
		log.Println("set remote description failed: ", err)
		return
	}

	select {}
}
