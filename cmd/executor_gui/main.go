package main

import (
	"context"

	"screen_recording/internal/app/executor_gui"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	app := executor_gui.NewApp("屏幕快照", 200, 100)
	app.Init()
	app.Start(ctx)
	app.End()
}

// package main
//
// import (
// 	"log"
//
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// )
//
// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Choice Widgets")
//
// 	check := widget.NewCheck("Optional", func(value bool) {
// 		log.Println("Check set to", value)
// 	})
// 	radio := widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(value string) {
// 		log.Println("Radio set to", value)
// 	})
// 	combo := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
// 		log.Println("Select set to", value)
// 	})
//
// 	myWindow.SetContent(container.NewVBox(check, radio, combo))
// 	myWindow.ShowAndRun()
// }
