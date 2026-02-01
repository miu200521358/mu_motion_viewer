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

	"github.com/miu200521358/mu_motion_viewer/pkg/adapter/mpresenter/messages"
	"github.com/miu200521358/mu_motion_viewer/pkg/usecase/minteractor"
)

const (
	motionViewerWindowIndex = 0
	motionViewerModelIndex  = 0
)

// motionViewerState はmu_motion_viewerの画面状態を保持する。
type motionViewerState struct {
	translator i18n.II18n
	logger     logging.ILogger
	userConfig config.IUserConfig

	usecase *minteractor.MotionViewerUsecase

	player               *widget.MotionPlayer
	modelPicker          *widget.FilePicker
	motionPicker         *widget.FilePicker
	saveModelButton      *widget.MPushButton
	saveSafeMotionButton *widget.MPushButton
	okBoneList           *ListBoxWidget
	okMorphList          *ListBoxWidget
	ngBoneList           *ListBoxWidget
	ngMorphList          *ListBoxWidget

	modelPath  string
	motionPath string
	modelData  *model.PmxModel
	motionData *motion.VmdMotion
}

// newMotionViewerState は画面状態を初期化する。
func newMotionViewerState(translator i18n.II18n, logger logging.ILogger, userConfig config.IUserConfig, viewerUsecase *minteractor.MotionViewerUsecase) *motionViewerState {
	if logger == nil {
		logger = logging.DefaultLogger()
	}
	return &motionViewerState{
		translator: translator,
		logger:     logger,
		userConfig: userConfig,
		usecase:    viewerUsecase,
	}
}

// applyInitialPaths は初期パスをウィジェットに反映する。
func (s *motionViewerState) applyInitialPaths(initialMotionPath string) {
	if s == nil {
		return
	}
	if s.modelPicker != nil && s.userConfig != nil {
		values, err := s.userConfig.GetStringSlice(config.UserConfigKeyPmxHistory)
		if err == nil && len(values) > 0 {
			s.modelPicker.SetPath(values[0])
		}
	}
	if s.motionPicker != nil && initialMotionPath != "" {
		s.motionPicker.SetPath(initialMotionPath)
	}
}

// handleModelPathChanged はモデルパス変更を処理する。
func (s *motionViewerState) handleModelPathChanged(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
	if s == nil {
		return
	}
	s.modelPath = path

	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), nil)
		s.modelData = nil
		if cw != nil {
			cw.SetModel(motionViewerWindowIndex, motionViewerModelIndex, nil)
		}
		s.updateCheckLists()
		return
	}
	result, err := s.usecase.LoadModel(rep, path)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), err)
		s.modelData = nil
		if cw != nil {
			cw.SetModel(motionViewerWindowIndex, motionViewerModelIndex, nil)
		}
		s.updateCheckLists()
		return
	}

	modelData := (*model.PmxModel)(nil)
	if result != nil {
		modelData = result.Model
	}
	s.modelData = modelData
	if cw != nil {
		cw.SetModel(motionViewerWindowIndex, motionViewerModelIndex, modelData)
	}
	s.updateCheckLists()
}

// handleMotionPathChanged はモーションパス変更を処理する。
func (s *motionViewerState) handleMotionPathChanged(cw *controller.ControlWindow, rep io_common.IFileReader, path string) {
	if s == nil {
		return
	}
	s.motionPath = path

	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), nil)
		s.motionData = nil
		if cw != nil {
			cw.SetMotion(motionViewerWindowIndex, motionViewerModelIndex, nil)
		}
		s.updatePlayerStateWithFrame(nil, 0)
		s.updateCheckLists()
		return
	}
	motionResult, err := s.usecase.LoadMotion(rep, path)
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, "読み込み失敗"), err)
		s.motionData = nil
		if cw != nil {
			cw.SetMotion(motionViewerWindowIndex, motionViewerModelIndex, nil)
		}
		s.updatePlayerStateWithFrame(nil, 0)
		s.updateCheckLists()
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
		cw.SetMotion(motionViewerWindowIndex, motionViewerModelIndex, motionData)
	}
	s.updatePlayerStateWithFrame(motionData, maxFrame)
	s.updateCheckLists()
}

// updatePlayerStateWithFrame は再生UIを反映する。
func (s *motionViewerState) updatePlayerStateWithFrame(motionData *motion.VmdMotion, maxFrame motion.Frame) {
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

// updateCheckLists はOK/NG一覧を更新する。
func (s *motionViewerState) updateCheckLists() {
	if s == nil {
		return
	}
	result, err := minteractor.CheckExists(s.modelData, s.motionData)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("OK/NG判定に失敗しました: %s", err.Error())
		}
		return
	}
	if s.okBoneList != nil {
		if err := s.okBoneList.SetItems(result.OkBones); err != nil {
			if s.logger != nil {
				s.logger.Error("OKボーン一覧の更新に失敗しました: %s", err.Error())
			}
		}
	}
	if s.okMorphList != nil {
		if err := s.okMorphList.SetItems(result.OkMorphs); err != nil {
			if s.logger != nil {
				s.logger.Error("OKモーフ一覧の更新に失敗しました: %s", err.Error())
			}
		}
	}
	if s.ngBoneList != nil {
		if err := s.ngBoneList.SetItems(result.NgBones); err != nil {
			if s.logger != nil {
				s.logger.Error("NGボーン一覧の更新に失敗しました: %s", err.Error())
			}
		}
	}
	if s.ngMorphList != nil {
		if err := s.ngMorphList.SetItems(result.NgMorphs); err != nil {
			if s.logger != nil {
				s.logger.Error("NGモーフ一覧の更新に失敗しました: %s", err.Error())
			}
		}
	}
}

// saveModelSetting は設定保存のログ出力のみを行う。
func (s *motionViewerState) saveModelSetting() {
	if s == nil {
		return
	}
	path := s.modelPath
	if s.usecase == nil || !s.usecase.CanLoadModelPath(path) {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.LogSaveFailure), nil)
		logInfoLine(s.logger, messages.LogSaveFailureDetail, path)
		controller.Beep()
		return
	}

	logInfoLine(s.logger, messages.LogSaveSuccess)
	logInfoLine(s.logger, messages.LogSaveSuccessDetail, path)
	controller.Beep()
}

// saveSafeMotion はIK無効モーションを保存する。
func (s *motionViewerState) saveSafeMotion() {
	if s == nil || s.motionData == nil {
		return
	}
	if s.usecase == nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.LogSafeSaveFailure), nil)
		controller.Beep()
		return
	}
	result, err := s.usecase.SaveSafeMotion(minteractor.SafeMotionSaveRequest{
		Motion:       s.motionData,
		FallbackPath: s.motionPath,
	})
	basePath := ""
	safePath := ""
	if result != nil {
		basePath = result.BasePath
		safePath = result.SafePath
	}
	if err != nil {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.LogSafeSaveFailure), err)
		logInfoLine(s.logger, messages.LogSafeSaveFailureDetail, safePath)
		controller.Beep()
		return
	}
	if basePath == "" || safePath == "" {
		logErrorWithTitle(s.logger, i18n.TranslateOrMark(s.translator, messages.LogSafeSaveFailure), nil)
		logInfoLine(s.logger, messages.LogSafeSaveFailureDetail, basePath)
		controller.Beep()
		return
	}

	logInfoLine(s.logger, messages.LogSafeSaveSuccess)
	logInfoLine(s.logger, messages.LogSafeSaveSuccessDetail, safePath)
	controller.Beep()
}
