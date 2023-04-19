package server

import (
	"fmt"
	"log"
	"net/http"
	"screen_recording/internal/channel"
	"screen_recording/internal/model"
	"screen_recording/internal/util"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	channelName := r.URL.Query().Get("channel")
	// 获取频道
	c := channel.Get(channelName)
	if c == nil {
		w.Write([]byte("频道不存在"))
		log.Printf("频道不存在 [%s]", channelName)
		return
	}
	
	// 绑定
	client, err := c.BindSubscriber(&channel.Client{
		Name: util.RemoteIp(r),
		Channel: make(chan string, 100),
	})
	if err != nil {
		w.Write([]byte(err.Error()))
		log.Printf("频道关注失败,err:%s", err.Error())
		return
	}
	
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	for data := range client.Channel {
		msg := &model.MessageEvent{
			Id: util.RandStr(11),
			Event: "message",
			Data: data,
		}
		// 写入数据
		fmt.Fprintf(w, "%s", msg)
		// 刷新到响应
		w.(http.Flusher).Flush()
	}
}