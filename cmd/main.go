//go:build windows
// +build windows

package main

import (
	"embed"
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/miu200521358/mu_motion_viewer/pkg/ui"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/interface/viewer"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
)

var env string

func init() {
	runtime.LockOSThread()

	// システム上の25%の論理プロセッサを使用する
	runtime.GOMAXPROCS(max(1, int(runtime.NumCPU()/4)))

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(widget.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

func main() {
	appConfig := mconfig.LoadAppConfig(appFiles)
	appConfig.Env = env
	mi18n.Initialize(appI18nFiles)

	mApp := app.NewMApp(appConfig)
	mApp.RunViewerToControlChannel()
	mApp.RunControlToViewerChannel()

	go func() {
		// 操作ウィンドウは別スレッドで起動
		controlWindow := controller.NewControlWindow(appConfig, mApp.ControlToViewerChannel(), ui.GetMenuItems, 1)
		mApp.SetControlWindow(controlWindow)

		controlWindow.InitTabWidget()
		ui.NewToolState(mApp, controlWindow)

		consoleView := widget.NewConsoleView(controlWindow.MainWindow, 256, 50)
		log.SetOutput(consoleView)

		mApp.RunController()
	}()

	viewerWindow := viewer.NewViewWindow(
		mApp.ViewerCount(), appConfig, mApp, mApp.ViewerToControlChannel(),
		fmt.Sprintf("%s %s", appConfig.Name, appConfig.Version), nil)
	viewerWindow.GetWindow().SetCloseCallback(func(w *glfw.Window) { mApp.AppState().SetClosed(true) })

	mApp.AddViewWindow(viewerWindow)

	mApp.Center()
	mApp.RunViewer()
}
