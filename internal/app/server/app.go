package server

import (
	"context"
	"log"
	"net/http"

	"screen_recording/internal/app"
	"screen_recording/internal/channel"
	"screen_recording/internal/util"
)

func init() {
	app.Set("server", &App{})
}

type App struct {
	Name    string
	Address string
	ChannelNum int
}

func (a *App) Init() error {
	a.Name = app.GlobalConfig.GetString("name")
	a.Address = app.GlobalConfig.GetString("address")
	a.ChannelNum = app.GlobalConfig.GetInt("ChannelNum")
	return nil
}

func (a *App) Start(ctx context.Context) {
	log.Printf("app [%s] Start", a.Name)

	// 创建频道频道
	go a.initChannel(ctx)

	// 注册路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/report", reportHandler)
	mux.HandleFunc("/sse", sseHandler)
	
	// 启动服务
	server := &http.Server{
		Addr:    a.Address,
		Handler: mux,
	}
	server.ListenAndServe()
}

func (a *App) End() {
	log.Printf("app [%s] End", a.Name)
}

// 初始化频道
func(a *App) initChannel(ctx context.Context){
	for i := 0; i < a.ChannelNum; i++ {
		chann := &channel.Channel{
			Publisher: make(chan string),
		}
		name := util.RandStr(10)
		channel.Set(name, chann)
		go chann.Start(ctx)
		log.Printf("频道ID:%s", name)
	}
}