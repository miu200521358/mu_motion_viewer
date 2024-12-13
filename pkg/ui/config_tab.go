package ui

import (
	"strings"
	"sync"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mi18n"
	"github.com/miu200521358/mlib_go/pkg/mutils/mlog"
	"github.com/miu200521358/walk/pkg/walk"
	"github.com/miu200521358/win"
)

func newConfigTab(controlWindow *controller.ControlWindow, toolState *ToolState) {
	toolState.ConfigTab = widget.NewMTabPage(mi18n.T("設定"))
	controlWindow.AddTabPage(toolState.ConfigTab.TabPage)

	toolState.ConfigTab.SetLayout(walk.NewVBoxLayout())

	composite, err := walk.NewComposite(toolState.ConfigTab)
	if err != nil {
		widget.RaiseError(err)
	}
	composite.SetLayout(walk.NewVBoxLayout())

	// ラベル
	label, err := walk.NewTextLabel(composite)
	if err != nil {
		widget.RaiseError(err)
	}
	label.SetText(mi18n.T("表示用モデル設定説明"))

	toolState.PmxPicker = widget.NewPmxReadFilePicker(
		controlWindow,
		composite,
		"pmx",
		"Pmxファイル",
		"Pmxファイルを選択してください",
		"Pmxファイルの使い方")
	toolState.PmxPicker.ChangePath(toolState.ModelPath)

	toolState.PmxPicker.SetOnPathChanged(func(path string) {
		if canLoad, err := toolState.PmxPicker.CanLoad(); !canLoad {
			if err != nil {
				mlog.ET(mi18n.T("読み込み失敗"), err.Error())
			}
			return
		}

		resultChan := make(chan loadPmxResult, 1)
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			var loadResult loadPmxResult
			rep := repository.NewPmxPmxJsonRepository()
			if data, err := rep.Load(path); err != nil {
				loadResult.model = nil
				loadResult.err = err
				resultChan <- loadResult
				return
			} else {
				loadResult.model = data.(*pmx.PmxModel)
				loadResult.err = nil
				resultChan <- loadResult
			}
		}()

		// 非同期で結果を受け取る
		go func() {
			wg.Wait()
			close(resultChan)

			result := <-resultChan

			if result.err != nil {
				mlog.ET(mi18n.T("読み込み失敗"), result.err.Error())
			} else if result.model == nil {
				toolState.ModelPath = ""
				toolState.Model = nil
			} else {
				// 強制更新用にハッシュ設定
				result.model.SetRandHash()

				toolState.ModelPath = path
				toolState.Model = result.model

				toolState.App.SetFuncGetModels(func() [][]*pmx.PmxModel {
					return [][]*pmx.PmxModel{{toolState.Model}}
				})

				toolState.App.SetFuncGetMotions(func() [][]*vmd.VmdMotion {
					return [][]*vmd.VmdMotion{{toolState.Motion}}
				})
			}
		}()
	})

	// NGボーン
	{
		walk.NewVSeparator(composite)

		// ラベル
		label, err := walk.NewTextLabel(composite)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("NG使用ボーン"))

		// NGボーン
		ngBoneEdit, err := walk.NewTextEditWithStyle(composite, win.WS_VISIBLE|win.WS_VSCROLL|win.ES_MULTILINE|win.ES_READONLY)
		if err != nil {
			widget.RaiseError(err)
		}
		ngBoneEdit.SetText(strings.Join(toolState.ActiveMissingBoneNames, "\r\n"))
	}

	// OKボーン
	{
		walk.NewVSeparator(composite)

		// ラベル
		label, err := walk.NewTextLabel(composite)
		if err != nil {
			widget.RaiseError(err)
		}
		label.SetText(mi18n.T("OK使用ボーン"))

		// OKボーン
		okBoneEdit, err := walk.NewTextEditWithStyle(composite, win.WS_VISIBLE|win.WS_VSCROLL|win.ES_MULTILINE|win.ES_READONLY)
		if err != nil {
			widget.RaiseError(err)
		}
		okBoneEdit.SetText(strings.Join(toolState.ActiveExistBoneNames, "\r\n"))
	}

	// フッター
	{
		walk.NewVSeparator(composite)

		playerComposite, err := walk.NewComposite(toolState.ConfigTab)
		if err != nil {
			widget.RaiseError(err)
		}
		playerComposite.SetLayout(walk.NewVBoxLayout())

		// プレイヤー
		toolState.Player = widget.NewMotionPlayer(playerComposite, controlWindow)
		toolState.Player.SetOnTriggerPlay(func(playing bool) { toolState.onPlay(playing) })
		controlWindow.SetPlayer(toolState.Player)
		toolState.ControlWindow.UpdateMaxFrame(toolState.Motion.MaxFrame())
	}

	{
		walk.NewVSeparator(composite)

		// 保存ボタン
		saveButton, err := walk.NewPushButton(composite)
		if err != nil {
			widget.RaiseError(err)
		}
		saveButton.SetText(mi18n.T("表示用モデル設定"))

		walk.NewVSpacer(toolState.ConfigTab)

		saveButton.Clicked().Attach(func() {
			if isOk, err := toolState.PmxPicker.CanLoad(); !isOk {
				if err != nil {
					mlog.ET(mi18n.T("保存失敗"), mi18n.T("保存エラーメッセージ",
						map[string]interface{}{"Error": err.Error()}))
				} else {
					mlog.ET(mi18n.T("保存失敗"), mi18n.T("保存失敗メッセージ",
						map[string]interface{}{"Path": toolState.PmxPicker.GetPath()}))
				}

				widget.Beep()

				return
			}

			widget.Beep()

			toolState.App.SetClosed(true)
		})
	}

	// 自動的にモーション再生
	toolState.Player.SetPlaying(true)
}
