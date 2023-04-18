package server

import (
	"context"
	"log"
	"net/http"

	"screen_recording/internal/app"
	"screen_recording/internal/model"
)

func init() {
	app.Set("server", &App{})
}

type App struct {
	Name    string
	Address string
	Channel map[string]chan *model.MessageEvent
}

func (a *App) Init() error {
	a.Name = app.GlobalConfig.GetString("name")
	a.Address = app.GlobalConfig.GetString("address")
	return nil
}

func (a *App) Start(ctx context.Context) {
	log.Printf("app [%s] Start", a.Name)

	// 注册路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewHandler)
	mux.HandleFunc("/report", reportHandler)

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
