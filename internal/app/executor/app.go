package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"flag"
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

var flagText = `
 _____                               _           _   
/ ____|                             | |         | |  
| (___   ___ _ __ ___  ___ _ __  ___| |__   ___ | |_ 
\___ \ /  __| '__/ _ \/ _ \ '_ \/ __| '_ \ / _ \| __|
____) | ( __| | |  __/  __/ | | \__ \ | | | (_) | |_ 
|_____/\ ___|_|  \___|\___|_| |_|___/_| |_|\___/ \__|																															
`

var logger = &log.Logger{}
func init() {
	logf, _ := os.Open("./log.log")
	logger.SetOutput(logf)
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
	a.Name = "屏幕快照"
	// 命令行参数覆盖
	flag.UintVar(&a.Frequency, "f", 0, "快照频率,单位秒")
	flag.UintVar(&a.Compression, "c", 100, "快照压缩比例(1-100)")
	flag.StringVar(&a.ReportChannel, "channel", "channel_1", "上报频道")
	flag.StringVar(&a.ReportAddress, "addr", "129.211.212.5:9999", "上报地址")
	flag.StringVar(&a.ReportKey, "key", "", "上报秘钥")
	flag.Parse()
	return nil
}

// Start 启动App
func (a *App) Start(ctx context.Context) {
	logger.Printf("app [%s] Start", a.Name)
	fmt.Println(flagText)
	fmt.Println("上报频道："+ a.ReportChannel)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Duration(a.Frequency) * time.Second)
			snapshot := a.snapshot()
			if err := a.report(snapshot); err != nil {
				logger.Printf("快照上报失败,err:%s", err.Error())
			}
		}
	}
}

// End 结束App
func (a *App) End() {
	logger.Printf("app [%s] End", a.Name)
}

// 屏幕快照
func (a *App) snapshot() []byte {
	// 获取第一个屏幕
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		logger.Printf("获取快照失败,err:%s", err.Error())
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
		logger.Printf("读取图片内容失败,err:%s", err.Error())
		return nil
	}
	return []byte(base64.StdEncoding.EncodeToString(srcByte))
}

// 上报快照信息
func (a *App) report(data []byte) error {
	// gzip压缩
	data, err := util.Gzip(data)
	if err != nil {
		return err
	}
	// 上报请求
	api := fmt.Sprintf("http://%s/report?channel=%s&key=%s", a.ReportAddress, a.ReportChannel, a.ReportKey)
	rsp, err := http.Post(api, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return nil
}
