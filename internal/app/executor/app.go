package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"screen_recording/internal/app"
	"screen_recording/internal/util"

	"time"

	"github.com/kbinani/screenshot"
)

func init() {
	app.Set("executor", &App{})
}

type App struct {
	Name          string
	Frequency     uint // 快照频率
	Compression   uint // 压缩比例，1-100
	ReportChannel string
	ReportAddress string
	ReportKey     string
}

// Init 初始化App
func (a *App) Init() error {
	a.Name = app.GlobalConfig.GetString("name")
	a.Frequency = app.GlobalConfig.GetUint("frequency")
	a.Compression = app.GlobalConfig.GetUint("compression")
	a.ReportChannel = app.GlobalConfig.GetString("ReportChannel")
	a.ReportAddress = app.GlobalConfig.GetString("ReportAddress")
	a.ReportKey = app.GlobalConfig.GetString("ReportKey")
	return nil
}

// Start 启动App
func (a *App) Start(ctx context.Context) {
	log.Printf("app [%s] Start", a.Name)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(a.Frequency) * time.Second)
			snapshot := a.snapshot()
			if err := a.report(snapshot); err != nil {
				log.Printf("快照上报失败,err:%s", err.Error())
			}
		}
	}
}

// End 结束App
func (a *App) End() {
	log.Printf("app [%s] End", a.Name)
}

// 屏幕快照
func (a *App) snapshot() []byte {
	// 获取第一个屏幕
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		log.Printf("获取快照失败，err:%s", err.Error())
		return nil
	}
	// 将图片写入临时文件
	fileName := fmt.Sprintf("window_%dx%d.jpeg", bounds.Dx(), bounds.Dy())
	file, _ := os.Create(fileName)
	defer file.Close()
	// 压缩比例100
	jpeg.Encode(file, img, &jpeg.Options{Quality: int(a.Compression)})

	// 转换为base64
	srcByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("读取图片内容失败，err:%s", err.Error())
		return nil
	}
	return []byte(base64.StdEncoding.EncodeToString(srcByte))
}

// 上报快照信息
func (a *App) report(data []byte) error {
	data, err := util.Gzip(data)
	if err != nil {
		return err
	}
	api := fmt.Sprintf("%s?channel=%s&key=%s", a.ReportAddress, a.ReportChannel, a.ReportKey)
	log.Println(api)
	log.Println(string(data))
	rsp, err := http.Post(api, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return nil
}
