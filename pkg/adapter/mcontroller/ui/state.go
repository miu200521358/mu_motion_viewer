//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/walk/pkg/walk"

	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mu_tree_viewer/pkg/usecase/minteractor"
)

const (
	treeViewerWindowIndex = 0
	treeViewerModelIndex  = 0
	folderHistoryKey      = "folder"
)

// treeViewerState はmu_tree_viewerの画面状態を保持する。
type treeViewerState struct {
	translator i18n.II18n
	logger     logging.ILogger
	userConfig config.IUserConfig

	usecase *minteractor.TreeViewerUsecase

	window       *controller.ControlWindow
	player       *widget.MotionPlayer
	folderPicker *FolderPicker
	motionPicker *widget.FilePicker
	treeView     *TreeViewWidget

	folderPaths []string
	motionPath  string
	modelData   *model.PmxModel
	motionData  *motion.VmdMotion
}

// newTreeViewerState は画面状態を初期化する。
func newTreeViewerState(translator i18n.II18n, logger logging.ILogger, userConfig config.IUserConfig, viewerUsecase *minteractor.TreeViewerUsecase) *treeViewerState {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	return &treeViewerState{
		translator: translator,
		logger:     logger,
		userConfig: userConfig,
		usecase:    viewerUsecase,
	}
}

// applyInitialPaths は初期パスをウィジェットに反映する。
func (s *treeViewerState) applyInitialPaths(initialMotionPath string) {
	if s == nil {
		return
	}
	if s.motionPicker != nil && initialMotionPath != "" {
		s.motionPicker.SetPath(initialMotionPath)
	}
}

// attachDropFiles はウィンドウのD&Dイベントを設定する。
func (s *treeViewerState) attachDropFiles(cw *controller.ControlWindow) {
	if s == nil || cw == nil {
		return
	}
	s.window = cw
	cw.DropFiles().Attach(func(files []string) {
		s.handleDropFiles(cw, files)
	})
}

// handleDropFiles はフォルダD&Dを処理する。
func (s *treeViewerState) handleDropFiles(cw *controller.ControlWindow, files []string) {
	if s == nil || len(files) == 0 {
		return
	}
	paths := cleanPaths(files)
	if len(paths) == 0 {
		return
	}
	if s.folderPicker != nil {
		s.folderPicker.SetPaths(paths)
		return
	}
	s.handleFolderPathsChanged(cw, paths)
}

// handleFolderPathsChanged はフォルダパス変更を処理する。
func (s *treeViewerState) handleFolderPathsChanged(cw *controller.ControlWindow, paths []string) {
	if s == nil {
		return
	}
	s.folderPaths = paths
	if s.treeView == nil {
		return
	}
	err := s.treeView.SetModelPaths(paths)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.LogTreeBuildFailure), err)
	}
	if len(paths) == 0 {
		return
	}
	if s.treeView.model == nil || s.treeView.model.RootCount() == 0 {
		logInfoLine(s.logger, i18n.TranslateOrMark(s.translator, messages.LogTreeEmpty))
	}
}

// handleMotionPathChanged はモーションパス変更を処理する。
func (s *treeViewerState) handleMotionPathChanged(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
	if s == nil {
		return
	}
	s.motionPath = path

	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.MessageLoadFailed), nil)
		s.motionData = nil
		if cw != nil {
			cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, nil)
		}
		s.updatePlayerStateWithFrame(nil, 0)
		return
	}
	motionResult, err := s.usecase.LoadMotion(rep, path)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.MessageLoadFailed), err)
		s.motionData = nil
		if cw != nil {
			cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, nil)
		}
		s.updatePlayerStateWithFrame(nil, 0)
		return
	}

	motionData := (*motion.VmdMotion)(nil)
	maxFrame := motion.Frame(0)
	if motionResult != nil {
		motionData = motionResult.Motion
		maxFrame = motionResult.MaxFrame
	}
	s.motionData = motionData
	if cw != nil {
		cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, motionData)
	}
	s.updatePlayerStateWithFrame(motionData, maxFrame)
}

// handleTreeFileSelected はツリーで選択されたモデルを読み込む。
func (s *treeViewerState) handleTreeFileSelected(path string) {
	if s == nil || path == "" {
		return
	}
	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.MessageLoadFailed), nil)
		return
	}
	result, err := s.usecase.LoadModel(nil, path)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.MessageLoadFailed), err)
		return
	}
	modelData := (*model.PmxModel)(nil)
	if result != nil {
		modelData = result.Model
	}
	s.modelData = modelData
	cw := s.controlWindow()
	if cw != nil {
		cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, modelData)
	}
	if modelData != nil {
		logInfoLine(s.logger, i18n.TranslateOrMark(s.translator, messages.LogLoadSuccess))
	}
}

// handleCopyPath はパスコピーを処理する。
func (s *treeViewerState) handleCopyPath(path string) {
	if s == nil || path == "" {
		return
	}
	if err := walk.Clipboard().SetText(path); err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.LogCopyFailure), err)
		return
	}
	logInfoLine(s.logger, i18n.TranslateOrMark(s.translator, messages.LogCopySuccess))
}

// updatePlayerStateWithFrame は再生UIを反映する。
func (s *treeViewerState) updatePlayerStateWithFrame(motionData *motion.VmdMotion, maxFrame motion.Frame) {
	if s == nil || s.player == nil {
		return
	}
	if motionData == nil {
		s.player.SetPlaying(false)
		s.player.Reset(0)
		return
	}
	if maxFrame <= 0 {
		maxFrame = motionData.MaxFrame()
	}
	s.player.Reset(maxFrame)
	if motionData.IsVpd() {
		s.player.SetPlaying(false)
		return
	}
	s.player.SetPlaying(true)
}

// controlWindow は関連付けられたウィンドウを取得する。
func (s *treeViewerState) controlWindow() *controller.ControlWindow {
	if s == nil {
		return nil
	}
	return s.window
}
