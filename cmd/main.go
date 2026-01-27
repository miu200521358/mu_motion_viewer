//go:build windows
// +build windows

package main

import (
	"embed"
	"os"
	"runtime"

	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_motion_viewer/pkg/ui"

	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/infra/app"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	sharedconfig "github.com/miu200521358/mlib_go/pkg/shared/base/config"
)

var env string

// init はOSスレッド固定とコンソール登録を行う。
func init() {
	runtime.LockOSThread()

	walk.AppendToWalkInit(func() {
		walk.MustRegisterWindowClass(controller.ConsoleViewClass)
	})
}

//go:embed app/*
var appFiles embed.FS

//go:embed i18n/*
var appI18nFiles embed.FS

// main はmu_motion_viewerを起動する。
func main() {
	initialMotionPath := app.FindInitialPath(os.Args, ".vmd", ".vpd")

	app.Run(app.RunOptions{
		ViewerCount: 1,
		AppFiles:    appFiles,
		I18nFiles:   appI18nFiles,
		AdjustConfig: func(appConfig *sharedconfig.AppConfig) {
			if env != "" {
				appConfig.EnvValue = sharedconfig.AppEnv(env)
			}
		},
		BuildMenuItems: func(baseServices base.IBaseServices) []declarative.MenuItem {
			return ui.NewMenuItems(baseServices.I18n(), baseServices.Logger())
		},
		BuildTabPages: func(widgets *controller.MWidgets, baseServices base.IBaseServices, audioPlayer audio_api.IAudioPlayer) []declarative.TabPage {
			return ui.NewTabPages(widgets, baseServices, initialMotionPath, audioPlayer)
		},
	})
}
