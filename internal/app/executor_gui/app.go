package executor_gui

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
	"github.com/spf13/cast"
	"screen_recording/internal/app"
	"screen_recording/internal/call"
)

var logger = &log.Logger{}

func init() {
	logf, _ := os.Create("./log.log")
	logger.SetOutput(logf)
	// logger.SetOutput(os.Stdout)

	os.Setenv("FYNE_THEME", "light")

	os.Setenv("FYNE_FONT", "C:\\Windows\\Fonts\\simhei.ttf")
}

type App struct {
	title      string
	w          float32
	h          float32
	app        fyne.App
	mainWindow fyne.Window

	serverAdd string
	channel   string
	frequency string
	quality   float64
}

func NewApp(title string, w, h float32) app.App {
	return &App{
		title:   title,
		w:       w,
		h:       h,
		quality: 0.4,
	}
}

func (a *App) Init() error {
	a.app = fyneapp.New()
	a.mainWindow = a.app.NewWindow(a.title)

	a.mainWindow.Resize(fyne.NewSize(a.w, a.h))
	a.mainWindow.SetFixedSize(true)
	return nil
}

func (a *App) Start(ctx context.Context) {
	logger.Println("app [屏幕快照] Start")
	// 服务器地址输入区域
	serverAddrEl := widget.NewEntry()
	serverAddrEl.SetText("127.0.0.1:9999")
	serverAddrArea := container.NewVBox(widget.NewLabel("服务器地址"), serverAddrEl)

	// 上报频道选择区域
	// channels := a.getChannels()
	channelSelectEl := widget.NewSelect([]string{
		"channel_1",
		"channel_2",
		"channel_3",
	}, func(value string) {
		// logger.Println("Select set to", value)
	})
	channelSelectEl.Selected = "channel_1"
	channelSelectArea := container.NewVBox(widget.NewLabel("选择上报频道"), channelSelectEl)

	// 上报频率选择区域
	frequencySelectEl := widget.NewRadioGroup([]string{"0.1", "0.4", "0.8", "1", "1.5"}, func(value string) {
		a.frequency = value
	})
	frequencySelectEl.Selected = "0.4"
	frequencySelectEl.Horizontal = true
	frequencySelectArea := container.NewVBox(widget.NewLabel("上报频率 (秒)"), frequencySelectEl)

	// 上报质量选择区域
	qualitySelectEl := widget.NewSlider(1, 100)
	qualitySelectEl.Step = 10
	qualitySelectEl.Value = 80
	qualityLabel := widget.NewLabel("快照质量(80%)")
	qualitySelectEl.OnChanged = func(f float64) {
		qualityLabel.SetText(fmt.Sprintf("快照质量(%.0f%%)", f))
		a.quality = f
	}
	qualitySelectArea := container.NewVBox(qualityLabel, qualitySelectEl)

	startReportBut := widget.NewButton("开始上报", nil)
	stopReportBut := widget.NewButton("停止上报", nil)
	stopReportBut.Hide()
	// 开始上报
	startReportBut.OnTapped = func() {
		a.serverAdd = serverAddrEl.Text
		a.channel = channelSelectEl.Selected
		logger.Println("服务器地址:", serverAddrEl.Text)
		logger.Println("上报频道:", channelSelectEl.Selected)
		logger.Println("上报频率 (秒):", frequencySelectEl.Selected)
		logger.Println("快照质量:", qualitySelectEl.Value)

		serverAddrEl.Disable()
		channelSelectEl.Disable()
		startReportBut.Hide()
		stopReportBut.Show()

		ctx, cancel := context.WithCancel(ctx)
		stopReportBut.OnTapped = func() {
			serverAddrEl.Enable()
			channelSelectEl.Enable()
			stopReportBut.Hide()
			startReportBut.Show()
			cancel()
		}
		go a.startReport(ctx)
	}

	content := container.NewVBox(
		serverAddrArea,
		channelSelectArea,
		frequencySelectArea,
		qualitySelectArea,
		startReportBut,
		stopReportBut,
	)
	a.mainWindow.SetContent(content)
	a.mainWindow.ShowAndRun()
}

func (a *App) End() {}

func (a *App) updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}

func (a *App) getChannels() []string {
	names, err := call.GetChannel("127.0.0.1")
	if err != nil {
		logger.Fatalln(err)
	}
	return names
}

func (a *App) startReport(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:

			time.Sleep(time.Duration(cast.ToFloat64(a.frequency) * float64(time.Second)))
			data, err := a.snapshot()
			if err != nil {
				logger.Printf("快照获取失败，err:%s", err.Error())
				continue
			}
			err = call.Report(data, a.serverAdd, a.channel)
			if err != nil {
				logger.Printf("快照获取失败，err:%s", err.Error())
				continue
			}
		}
	}
}

func (a *App) snapshot() (io.Reader, error) {
	// 获取第一个屏幕
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	// 压缩比例
	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: int(a.quality)})
	if err != nil {
		return nil, err
	}
	return buf, nil
}
