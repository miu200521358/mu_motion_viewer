// 指示: miu200521358
package minteractor

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mu_motion_viewer/pkg/usecase/port/moutput"
)

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

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep moutput.IFileReader, path string) (*model.PmxModel, error) {
	return commonusecase.LoadModel(rep, path)
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep moutput.IFileReader, path string) (*motion.VmdMotion, error) {
	return commonusecase.LoadMotion(rep, path)
}

// LoadModelWithValidation はモデルを読み込み、結果を返す。
// テクスチャ検証は mu_model_viewer 固有処理のため、ここでは実行しない。
func LoadModelWithValidation(rep moutput.IFileReader, path string, validator moutput.ITextureValidator) (*ModelLoadResult, error) {
	_ = validator // テクスチャ検証はmu_model_viewer固有のため、このusecaseでは扱わない。
	modelData, err := LoadModel(rep, path)
	if err != nil {
		return nil, err
	}
	return &ModelLoadResult{Model: modelData}, nil
}

// LoadMotionWithMeta はモーションを読み込み、最大フレームを返す。
func LoadMotionWithMeta(rep moutput.IFileReader, path string) (*MotionLoadResult, error) {
	return commonusecase.LoadMotionWithMeta(rep, path)
}

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func CanLoadPath(rep moutput.IFileReader, path string) bool {
	return commonusecase.CanLoadPath(rep, path)
}
