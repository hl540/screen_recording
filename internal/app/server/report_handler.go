package server

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"screen_recording/internal/channel"
)

func reportHandler(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Query().Get("channel")
	c := channel.Get(channelName)
	if c == nil {
		w.Write([]byte("频道不存在"))
		log.Printf("频道不存在 [%s]", channelName)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Printf("上报内容解析失败,err:%s", err.Error())
		return
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Printf("上报内容解析失败,err:%s", err.Error())
		return
	}
	c.Publisher <- base64.StdEncoding.EncodeToString(data)
	w.Write([]byte("success"))
}
