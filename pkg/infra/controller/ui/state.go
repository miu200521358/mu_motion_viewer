//go:build windows
// +build windows

// 指示: miu200521358
package ui

import (
	"path/filepath"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	"github.com/miu200521358/mlib_go/pkg/infra/controller"
	"github.com/miu200521358/mlib_go/pkg/infra/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/shared/base/config"
	"github.com/miu200521358/mlib_go/pkg/shared/base/i18n"
	"github.com/miu200521358/mlib_go/pkg/shared/base/logging"
	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mu_tree_viewer/pkg/adapter/mtree"
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

	folderPicker *FolderPicker
	motionPicker *widget.FilePicker
	treeView     *TreeViewWidget
	player       *widget.MotionPlayer

	pendingSelectPath string
	modelPath         string
	motionPath        string
	modelData         *model.PmxModel
	motionData        *motion.VmdMotion
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
func (s *treeViewerState) applyInitialPaths(initialModelPath string, initialMotionPath string) {
	if s == nil {
		return
	}
	if initialModelPath != "" {
		s.pendingSelectPath = initialModelPath
		dir := filepath.Dir(initialModelPath)
		if s.folderPicker != nil {
			s.folderPicker.SetPaths([]string{dir})
		}
	} else if s.folderPicker != nil && s.userConfig != nil {
		values, err := s.userConfig.GetStringSlice(folderHistoryKey)
		if err == nil && len(values) > 0 {
			s.folderPicker.SetPaths([]string{values[0]})
		}
	}
	if s.motionPicker != nil && initialMotionPath != "" {
		s.motionPicker.SetPath(initialMotionPath)
	}
}

// handleFolderPathsChanged はフォルダパス変更を処理する。
func (s *treeViewerState) handleFolderPathsChanged(cw *controller.ControlWindow, paths []string) {
	if s == nil {
		return
	}
	roots, err := mtree.BuildModelTree(paths)
	if err != nil && s.logger != nil {
		s.logger.Error("フォルダツリーの構築に失敗しました: %s", err.Error())
	}
	if s.treeView != nil {
		s.treeView.SetRoots(roots)
	}
	if cw != nil {
		cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, nil)
	}
	if s.pendingSelectPath != "" && s.treeView != nil {
		if s.treeView.SelectPath(s.pendingSelectPath) {
			s.pendingSelectPath = ""
		}
	}
}

// handleFoldersDropped はフォルダのD&Dを処理する。
func (s *treeViewerState) handleFoldersDropped(_ *controller.ControlWindow, files []string) {
	if s == nil || s.folderPicker == nil {
		return
	}
	s.folderPicker.SetPaths(files)
}

// handleFileSelected はツリー上のファイル選択を処理する。
func (s *treeViewerState) handleFileSelected(cw *controller.ControlWindow, path string) {
	if s == nil {
		return
	}
	s.modelPath = path
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, nil)
		return
	}
	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), nil)
		cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, nil)
		return
	}
	result, err := s.usecase.LoadModel(nil, path)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), err)
		cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, nil)
		return
	}
	modelData := (*model.PmxModel)(nil)
	if result != nil {
		modelData = result.Model
	}
	s.modelData = modelData
	if modelData == nil {
		cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, nil)
		return
	}
	cw.SetModel(treeViewerWindowIndex, treeViewerModelIndex, modelData)
	logInfoLine(s.logger, i18n.TranslateOrMark(s.translator, messages.LogLoadSuccess))
}

// handleMotionPathChanged はモーションパス変更を処理する。
func (s *treeViewerState) handleMotionPathChanged(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
	if s == nil {
		return
	}
	s.motionPath = path
	if cw == nil {
		return
	}
	if path == "" {
		cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, nil)
		s.updatePlayerStateWithFrame(nil, 0)
		return
	}
	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), nil)
		cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, nil)
		s.updatePlayerStateWithFrame(nil, 0)
		return
	}
	motionResult, err := s.usecase.LoadMotion(rep, path)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), err)
		cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, nil)
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
	if motionData == nil {
		cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, nil)
		s.updatePlayerStateWithFrame(nil, 0)
		return
	}
	cw.SetMotion(treeViewerWindowIndex, treeViewerModelIndex, motionData)
	s.updatePlayerStateWithFrame(motionData, maxFrame)
	logInfoLine(s.logger, i18n.TranslateOrMark(s.translator, messages.LogLoadSuccess))
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
