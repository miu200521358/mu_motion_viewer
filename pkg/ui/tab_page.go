package ui

import (
	"github.com/miu200521358/mlib_go/pkg/config/mi18n"
	"github.com/miu200521358/mlib_go/pkg/config/mlog"
	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mu_motion_viewer/pkg/usecase"
	"github.com/miu200521358/walk/pkg/declarative"
	"github.com/miu200521358/walk/pkg/walk"
)

func NewTabPage(mWidgets *controller.MWidgets) declarative.TabPage {
	var fileTab *walk.TabPage

	player := widget.NewMotionPlayer()

	var okBoneNamesListbox *walk.ListBox
	var okMorphNamesListbox *walk.ListBox
	var ngBoneNamesListbox *walk.ListBox
	var ngMorphNamesListbox *walk.ListBox

	pmxLoadPicker := widget.NewPmxXLoadFilePicker(
		"pmx",
		mi18n.T("モデルファイル"),
		mi18n.T("モデルファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if path == "" {
				cw.StoreModel(0, 0, nil)
				return
			}

			if data, err := rep.Load(path); err == nil {
				model := data.(*pmx.PmxModel)
				cw.StoreModel(0, 0, model)

				motion := cw.LoadMotion(0, 0)
				okBoneNames, okMorphNames, ngBoneNames, ngMorphNames := usecase.CheckExists(model, motion)
				okBoneNamesListbox.SetModel(okBoneNames)
				okMorphNamesListbox.SetModel(okMorphNames)
				ngBoneNamesListbox.SetModel(ngBoneNames)
				ngMorphNamesListbox.SetModel(ngMorphNames)
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err, "")
			}
		},
	)

	vmdLoadPicker := widget.NewVmdVpdLoadFilePicker(
		"vmd",
		mi18n.T("モーションファイル"),
		mi18n.T("モーションファイルを選択してください"),
		func(cw *controller.ControlWindow, rep repository.IRepository, path string) {
			if path == "" {
				cw.StoreMotion(0, 0, nil)
				return
			}

			if data, err := rep.Load(path); err == nil {
				motion := data.(*vmd.VmdMotion)
				player.Reset(motion.MaxFrame())
				cw.StoreMotion(0, 0, motion)

				model := cw.LoadModel(0, 0)
				okBoneNames, okMorphNames, ngBoneNames, ngMorphNames := usecase.CheckExists(model, motion)
				okBoneNamesListbox.SetModel(okBoneNames)
				okMorphNamesListbox.SetModel(okMorphNames)
				ngBoneNamesListbox.SetModel(ngBoneNames)
				ngMorphNamesListbox.SetModel(ngMorphNames)

				if !motion.IsVpd() {
					// ポースじゃなければ再生
					player.SetPlaying(true)
				}
				// フォーカスを当てる
				cw.SetFocus()
			} else {
				mlog.ET(mi18n.T("読み込み失敗"), err, "")
			}
		},
	)

	var saveButton *walk.PushButton

	mWidgets.Widgets = append(mWidgets.Widgets, player, pmxLoadPicker, vmdLoadPicker)
	mWidgets.SetOnLoaded(func() {
		// 読み込みが完了したら、モデルのパスを設定
		if modelPath, motionPath, err := usecase.LoadModelMotionPath(); err == nil {
			pmxLoadPicker.SetPath(modelPath)
			vmdLoadPicker.SetPath(motionPath)
		}
	})

	player.SetOnChangePlayingPre(func(playing bool) {
		// ボーン有無などはいつでも操作OK
		okBoneNamesListbox.SetEnabled(true)
		okMorphNamesListbox.SetEnabled(true)
		ngBoneNamesListbox.SetEnabled(true)
		ngMorphNamesListbox.SetEnabled(true)
		saveButton.SetEnabled(playing)
	})

	return declarative.TabPage{
		Title:    mi18n.T("ファイル"),
		AssignTo: &fileTab,
		Layout:   declarative.VBox{},
		Background: declarative.SolidColorBrush{
			Color: controller.ColorTabBackground,
		},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.VBox{},
				Children: []declarative.Widget{
					pmxLoadPicker.Widgets(),
					vmdLoadPicker.Widgets(),
					declarative.VSeparator{},
					declarative.Composite{
						Layout: declarative.Grid{
							Columns: 2,
						},
						Children: []declarative.Widget{
							declarative.Label{
								Text:        mi18n.T("OKボーン"),
								ToolTipText: mi18n.T("OKボーン説明"),
							},
							declarative.Label{
								Text:        mi18n.T("OKモーフ"),
								ToolTipText: mi18n.T("OKモーフ説明"),
							},
							declarative.ListBox{
								AssignTo: &okBoneNamesListbox,
							},
							declarative.ListBox{
								AssignTo: &okMorphNamesListbox,
							},
							declarative.Label{
								Text:        mi18n.T("NGボーン"),
								ToolTipText: mi18n.T("NGボーン説明"),
							},
							declarative.Label{
								Text:        mi18n.T("NGモーフ"),
								ToolTipText: mi18n.T("NGモーフ説明"),
							},
							declarative.ListBox{
								AssignTo: &ngBoneNamesListbox,
							},
							declarative.ListBox{
								AssignTo: &ngMorphNamesListbox,
							},
						},
					},
					declarative.VSpacer{},
					declarative.PushButton{
						AssignTo: &saveButton,
						Text:     mi18n.T("設定保存"),
						OnClicked: func() {
							if isOk := pmxLoadPicker.CanLoad(); !isOk {
								mlog.ET(mi18n.T("保存失敗"), nil, mi18n.T("保存失敗メッセージ",
									map[string]any{"Path": pmxLoadPicker.Path()}))
							} else {
								mlog.IT(mi18n.T("保存成功"), mi18n.T("保存成功メッセージ",
									map[string]any{"Path": pmxLoadPicker.Path()}))
							}

							controller.Beep()
						},
					},
					declarative.VSpacer{},
					player.Widgets(),
					declarative.VSpacer{},
				},
			},
		},
	}
}
