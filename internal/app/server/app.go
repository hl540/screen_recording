package server

import (
	"context"
	"log"
	"net/http"

	"screen_recording/internal/app"
	"screen_recording/internal/channel"

	"github.com/spf13/viper"
)

func init() {
	app.Set("server", &App{})
}

type App struct {
	Name     string
	Address  string
	Channels []string
}

func (a *App) Init() error {
	viper.SetConfigFile("./conf.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	a.Name = viper.GetString("name")
	a.Address = viper.GetString("address")
	a.Channels = viper.GetStringSlice("Channels")
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
	mux.HandleFunc("/channel/all", getChannelHandler)
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
func (a *App) initChannel(ctx context.Context) {
	for _, name := range a.Channels {
		chann := &channel.Channel{
			Publisher: make(chan string),
		}
		channel.Set(name, chann)
		go chann.Start(ctx)
		log.Printf("频道ID:%s", name)
	}
}
