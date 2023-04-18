package main

import (
	"screen_recording/internal/app"
	_ "screen_recording/internal/app/server"
)

func main() {
	app.Run()
}
