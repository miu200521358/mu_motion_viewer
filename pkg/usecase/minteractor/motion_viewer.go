// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mu_motion_viewer/pkg/usecase/port/moutput"
)

// ModelLoadResult はモデル読み込み結果を表す。
type ModelLoadResult = commonusecase.ModelLoadResult

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult = commonusecase.MotionLoadResult

// MotionViewerUsecaseDeps はモーションビューア用ユースケースの依存を表す。
type MotionViewerUsecaseDeps struct {
	ModelReader  moutput.IFileReader
	MotionReader moutput.IFileReader
	MotionWriter moutput.IFileWriter
}

// MotionViewerUsecase はモーションビューアの入出力処理をまとめたユースケースを表す。
type MotionViewerUsecase struct {
	modelReader  moutput.IFileReader
	motionReader moutput.IFileReader
	motionWriter moutput.IFileWriter
}

// NewMotionViewerUsecase はモーションビューア用ユースケースを生成する。
func NewMotionViewerUsecase(deps MotionViewerUsecaseDeps) *MotionViewerUsecase {
	return &MotionViewerUsecase{
		modelReader:  deps.ModelReader,
		motionReader: deps.MotionReader,
		motionWriter: deps.MotionWriter,
	}
}

// LoadModel はモデルを読み込み、結果を返す。
func (uc *MotionViewerUsecase) LoadModel(rep moutput.IFileReader, path string) (*ModelLoadResult, error) {
	repo := rep
	if repo == nil {
		repo = uc.modelReader
	}
	modelData, err := commonusecase.LoadModel(repo, path)
	if err != nil {
		return nil, err
	}
	return &ModelLoadResult{Model: modelData}, nil
}

// LoadMotion はモーションを読み込み、最大フレーム情報を返す。
func (uc *MotionViewerUsecase) LoadMotion(rep moutput.IFileReader, path string) (*MotionLoadResult, error) {
	repo := rep
	if repo == nil {
		repo = uc.motionReader
	}
	return commonusecase.LoadMotionWithMeta(repo, path)
}

// CanLoadModelPath はモデルの読み込み可否を判定する。
func (uc *MotionViewerUsecase) CanLoadModelPath(path string) bool {
	return commonusecase.CanLoadPath(uc.modelReader, path)
}

// SaveSafeMotion はIK無効モーションを保存する。
func (uc *MotionViewerUsecase) SaveSafeMotion(request SafeMotionSaveRequest) (*SafeMotionSaveResult, error) {
	if request.Writer == nil {
		request.Writer = uc.motionWriter
	}
	return SaveSafeMotion(request)
}

// ExtractModelData は読み込み結果からモデルを取り出す。
func ExtractModelData(result *ModelLoadResult) *model.PmxModel {
	if result == nil {
		return nil
	}
	return result.Model
}

// ExtractMotionData は読み込み結果からモーションと最大フレームを取り出す。
func ExtractMotionData(result *MotionLoadResult) (*motion.VmdMotion, motion.Frame) {
	if result == nil {
		return nil, 0
	}
	return result.Motion, result.MaxFrame
}
