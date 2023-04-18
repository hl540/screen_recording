package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App interface {
	Init() error
	Start(ctx context.Context)
	End()
}

type helloApp struct{}

func (a *helloApp) Init() error {
	return nil
}

func (a *helloApp) Start(ctx context.Context) {
	log.Printf("app [%s] Start", "hello")
	a.End()
}

func (a *helloApp) End() {
	log.Printf("app [%s] End", "hello")
}

var appMap = make(map[string]App)
var mux sync.Mutex

func Set(name string, app App) {
	mux.Lock()
	defer mux.Unlock()
	appMap[name] = app
}

func Get(name string) App {
	app, ok := appMap[name]
	if !ok {
		return &helloApp{}
	}
	return app
}

func Init() map[string]App {
	ready := make(map[string]App)
	for name, app := range appMap {
		err := app.Init()
		if err != nil {
			log.Printf("app [%s] Init err:%s", name, err)
		}
		ready[name] = app
	}
	return ready
}

func Run() {
	ctx, cancel := context.WithCancel(context.TODO())
	ready := Init()
	for name, app := range ready {
		go func(name string, app App) {
			defer func(name string) {
				if err := recover(); err != nil {
					log.Printf("app [%s] Start err:%s", name, err)
				}
			}(name)
			app.Start(ctx)
		}(name, app)
	}
	// 退出信号量
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-signalChan

	// 关闭上下文
	cancel()

	// 结束app
	for _, app := range ready {
		app.End()
	}
}
