package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
)

type MessageEvent struct {
	Id   string
	Name string
	Data string
}

func (e MessageEvent) String() string {
	f := "id:%s\n"
	f += "event:%s\n"
	f += "data:%s\n\n"
	return fmt.Sprintf(f, e.Id, e.Name, e.Data)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// 流管道
	var streamChan = make(chan MessageEvent)

	// 启动截屏服务
	go SnapshotService(ctx, streamChan)
	// 视图服务
	go ViewService(ctx, streamChan)

	// 退出信号量
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	// 关闭上下文
	cancel()
}

// 截屏服务
func SnapshotService(ctx context.Context, streamChan chan<- MessageEvent) {
	timeTicker := time.NewTicker(1)
	for {
		select {
		case <-timeTicker.C:
			// 截屏
			streamChan <- MessageEvent{
				Id:   strconv.FormatInt(time.Now().Unix(), 10),
				Name: "message",
				Data: Snapshot(),
			}
		case <-ctx.Done():
			return
		default:
			// 截屏
			streamChan <- MessageEvent{
				Id:   strconv.FormatInt(time.Now().Unix(), 10),
				Name: "message",
				Data: Snapshot(),
			}
		}
	}
}

// 屏幕快照
func Snapshot() string {
	// 获取第一个屏幕
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Printf("获取快照失败，err:%s", err.Error())
		return ""
	}
	// 将图片写入临时文件
	fileName := fmt.Sprintf("window_%dx%d.jpeg", bounds.Dx(), bounds.Dy())
	file, _ := os.Create(fileName)
	defer file.Close()
	// 压缩比例100
	jpeg.Encode(file, img, &jpeg.Options{Quality: 100})

	// gzip压缩转换为base64
	srcByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("读取图片内容失败，err:%s", err.Error())
		return ""
	}
	// base64Byte, err := GzipStream(srcByte)
	// if err != nil {
	// 	log.Printf("gzip压缩，err:%s", err.Error())
	// 	return ""
	// }
	return base64.StdEncoding.EncodeToString(srcByte)
}

// 视图服务
func ViewService(ctx context.Context, streamChan <-chan MessageEvent) {
	http.HandleFunc("/window", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./view.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		tmpl.Execute(w, nil)
	})
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		for stream := range streamChan {
			log.Println(len(stream.Data))
			// 写入数据
			fmt.Fprintf(w, "%s", &stream)
			// 刷新到响应
			w.(http.Flusher).Flush()
		}
	})
	http.ListenAndServe(":9999", nil)
}

func GzipStream(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	base64Str := base64.StdEncoding.EncodeToString(b.Bytes())
	return []byte(base64Str), nil
}
