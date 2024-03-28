package server

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"

	"screen_recording/internal/channel"
)

func reportHandler(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Query().Get("channel")
	msgChannel := channel.Get(channelName)
	if msgChannel == nil {
		w.Write([]byte("频道不存在"))
		log.Printf("频道不存在 [%s]", channelName)
		return
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Printf("上报内容解析失败[FormFile],err:%s", err.Error())
		return
	}
	data, err := io.ReadAll(file)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Printf("上报内容解析失败[ReadAll],err:%s", err.Error())
		return
	}
	msgChannel.Publisher <- base64.StdEncoding.EncodeToString(data)
	w.Write([]byte("success"))
}
