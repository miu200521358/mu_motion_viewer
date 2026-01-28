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

	"github.com/miu200521358/mu_motion_viewer/pkg/ui_messages_labels"
	"github.com/miu200521358/mu_motion_viewer/pkg/usecase"
)

// NewTabPages はmu_motion_viewer用のタブページを生成する。
func NewTabPages(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialMotionPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *usecase.MotionViewerUsecase) []declarative.TabPage {
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

	state := newMotionViewerState(translator, logger, userConfig, viewerUsecase)

	state.player = widget.NewMotionPlayer(translator)
	state.player.SetAudioPlayer(audioPlayer, userConfig)

	state.modelPicker = widget.NewPmxPmdXLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyPmxHistory,
		i18n.TranslateOrMark(translator, ui_messages_labels.LabelModelFile),
		i18n.TranslateOrMark(translator, ui_messages_labels.LabelModelFileTip),
		state.handleModelPathChanged,
	)
	state.motionPicker = widget.NewVmdVpdLoadFilePicker(
		userConfig,
		translator,
		config.UserConfigKeyVmdHistory,
		i18n.TranslateOrMark(translator, ui_messages_labels.LabelMotionFile),
		i18n.TranslateOrMark(translator, ui_messages_labels.LabelMotionFileTip),
		state.handleMotionPathChanged,
	)

	state.saveModelButton = widget.NewMPushButton()
	state.saveModelButton.SetLabel(i18n.TranslateOrMark(translator, ui_messages_labels.LabelSettingSave))
	state.saveModelButton.SetTooltip(i18n.TranslateOrMark(translator, ui_messages_labels.LabelSettingSave))
	state.saveModelButton.SetOnClicked(func(_ *controller.ControlWindow) {
		state.saveModelSetting()
	})

	state.saveSafeMotionButton = widget.NewMPushButton()
	state.saveSafeMotionButton.SetLabel(i18n.TranslateOrMark(translator, ui_messages_labels.LabelSafeMotionSave))
	state.saveSafeMotionButton.SetTooltip(i18n.TranslateOrMark(translator, ui_messages_labels.LabelSafeMotionSave))
	state.saveSafeMotionButton.SetOnClicked(func(_ *controller.ControlWindow) {
		state.saveSafeMotion()
	})

	listMinSize := declarative.Size{Width: 220, Height: 80}
	state.okBoneList = NewListBoxWidget(i18n.TranslateOrMark(translator, ui_messages_labels.LabelOkBoneTip), logger)
	state.okBoneList.SetMinSize(listMinSize)
	state.okBoneList.SetStretchFactor(1)

	state.okMorphList = NewListBoxWidget(i18n.TranslateOrMark(translator, ui_messages_labels.LabelOkMorphTip), logger)
	state.okMorphList.SetMinSize(listMinSize)
	state.okMorphList.SetStretchFactor(1)

	state.ngBoneList = NewListBoxWidget(i18n.TranslateOrMark(translator, ui_messages_labels.LabelNgBoneTip), logger)
	state.ngBoneList.SetMinSize(listMinSize)
	state.ngBoneList.SetStretchFactor(1)

	state.ngMorphList = NewListBoxWidget(i18n.TranslateOrMark(translator, ui_messages_labels.LabelNgMorphTip), logger)
	state.ngMorphList.SetMinSize(listMinSize)
	state.ngMorphList.SetStretchFactor(1)

	if mWidgets != nil {
		mWidgets.Widgets = append(mWidgets.Widgets,
			state.player,
			state.modelPicker,
			state.motionPicker,
			state.saveModelButton,
			state.saveSafeMotionButton,
			state.okBoneList,
			state.okMorphList,
			state.ngBoneList,
			state.ngMorphList,
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
			state.applyInitialPaths(initialMotionPath)
		})
	}

	fileTabPage := declarative.TabPage{
		Title:    i18n.TranslateOrMark(translator, ui_messages_labels.LabelFile),
		AssignTo: &fileTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					state.modelPicker.Widgets(),
					state.motionPicker.Widgets(),
				},
			},
			declarative.VSeparator{},
			declarative.Composite{
				Layout: declarative.Grid{
					Columns: 2,
				},
				Children: []declarative.Widget{
					buildListBoxColumn(
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelOkBone),
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelOkBoneTip),
						state.okBoneList,
					),
					buildListBoxColumn(
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelOkMorph),
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelOkMorphTip),
						state.okMorphList,
					),
					buildListBoxColumn(
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelNgBone),
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelNgBoneTip),
						state.ngBoneList,
					),
					buildListBoxColumn(
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelNgMorph),
						i18n.TranslateOrMark(translator, ui_messages_labels.LabelNgMorphTip),
						state.ngMorphList,
					),
				},
			},
			declarative.VSeparator{},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					state.saveModelButton.Widgets(),
					state.saveSafeMotionButton.Widgets(),
				},
			},
			declarative.VSeparator{},
			state.player.Widgets(),
			declarative.VSpacer{},
		},
	}

	return []declarative.TabPage{fileTabPage}
}

// NewTabPage はmu_motion_viewer用の単一タブを生成する。
func NewTabPage(mWidgets *controller.MWidgets, baseServices base.IBaseServices, initialMotionPath string, audioPlayer audio_api.IAudioPlayer, viewerUsecase *usecase.MotionViewerUsecase) declarative.TabPage {
	return NewTabPages(mWidgets, baseServices, initialMotionPath, audioPlayer, viewerUsecase)[0]
}

// buildListBoxColumn はラベル付きのリスト表示を構成する。
func buildListBoxColumn(label string, tooltip string, listBox *ListBoxWidget) declarative.Composite {
	return declarative.Composite{
		Layout: declarative.VBox{
			MarginsZero: true,
			SpacingZero: true,
		},
		StretchFactor: 1,
		Children: []declarative.Widget{
			declarative.TextLabel{
				Text:        label,
				ToolTipText: tooltip,
			},
			listBox.Widgets(),
		},
	}
}
