//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/audio_api"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/shared/base"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mu_tree_viewer/pkg/usecase/minteractor"
)

const (
	treeViewFixedHeight = 200
)

// NewTabPages はmu_tree_viewer用のタブページを生成する。
func NewTabPages(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialMotionPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *minteractor.TreeViewerUsecase) []declarative.TabPage {
	var fileTab *walk.TabPage

	var translator i18n.II18n
	var logger logging.ILogger
	var userConfig config.IUserConfig
	if baseServices != nil {
		translator = baseServices.I18n()
		logger = baseServices.Logger()
		if cfg := baseServices.Config(); cfg != nil {
			userConfig = cfg.UserConfig()
		}
	}
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	if viewerUsecase == nil {
		viewerUsecase = minteractor.NewTreeViewerUsecase(minteractor.TreeViewerUsecaseDeps{})
	}

	state := newTreeViewerState(translator, logger, userConfig, viewerUsecase)

	state.player = widget.NewMotionPlayer(translator)
	state.player.SetAudioPlayer(audioPlayer, userConfig)

	state.folderPicker = NewFolderPicker(
		userConfig,
		translator,
		folderHistoryKey,
		i18n.TranslateOrMark(translator, messages.LabelFolderPath),
		i18n.TranslateOrMark(translator, messages.LabelFolderPathTip),
		state.handleFolderPathsChanged,
	)

	state.motionPicker = widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyVmdHistory,
		i18n.TranslateOrMark(translator, messages.LabelMotionPath),
		i18n.TranslateOrMark(translator, messages.LabelMotionPathTip),
		state.handleMotionPathChanged,
	)

	state.treeView = NewTreeViewWidget(translator, logger, state.handleTreeFileSelected, state.handleCopyPath)
	state.treeView.SetMinSize(declarative.Size{Width: 400, Height: treeViewFixedHeight})
	state.treeView.SetStretchFactor(1)

	if mWidgets != nil {
		mWidgets.Widgets = append(mWidgets.Widgets,
			state.folderPicker,
			state.motionPicker,
			state.treeView,
			state.player,
		)
		mWidgets.SetOnLoaded(func() {
			if mWidgets == nil || mWidgets.Window() == nil {
				return
			}
			mWidgets.Window().SetOnEnabledInPlaying(func(playing bool) {
				for _, w := range mWidgets.Widgets {
					w.SetEnabledInPlaying(playing)
				}
			})
			state.attachDropFiles(mWidgets.Window())
			state.applyInitialPaths(initialMotionPath)
		})
	}

	fileTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, messages.LabelFile),
		AssignTo: &fileTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					state.folderPicker.Widgets(),
					state.motionPicker.Widgets(),
					declarative.VSeparator{},
					declarative.Composite{
						Layout: declarative.VBox{},
						Children: []declarative.Widget{
							declarative.TextLabel{Text: i18n.TranslateOrMark(translator, messages.LabelTreeView)},
							state.treeView.Widgets(),
						},
					},
					declarative.VSeparator{},
					state.player.Widgets(),
				},
			},
		},
	}

	return []declarative.TabPage{fileTabPage}
}

// NewTabPage はmu_tree_viewer用の単一タブを生成する。
func NewTabPage(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialMotionPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *minteractor.TreeViewerUsecase) declarative.TabPage {
	return NewTabPages(mWidgets, baseServices, initialMotionPath, audioPlayer, viewerUsecase)[0]
}
