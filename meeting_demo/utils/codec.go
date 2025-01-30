package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

// Encode 将 webrtc.SessionDescription 对象 转为 Base64 编码的 JSON 字符串
func Encode(obj interface{}) string {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Println("json marshal error:", err)
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// Decode 从 Base64 编码的 JSON 字符串还原回原始 webrtc.SessionDescription 对象
func Decode(inputStr string, obj interface{}) {
	bytes, err := base64.StdEncoding.DecodeString(inputStr)
	if err != nil {
		log.Println("base64 decode error: ", err)
		panic(err)
	}
	err = json.Unmarshal(bytes, obj)
	if err != nil {
		log.Println("json marshal error:", err)
		panic(err)
	}
}
