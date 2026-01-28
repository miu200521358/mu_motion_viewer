// 指示: miu200521358
package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	portio "github.com/miu200521358/mu_motion_viewer/pkg/usecase/port/io"
)

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep portio.IFileReader, path string) (*model.PmxModel, error) {
	return commonusecase.LoadModel(rep, path)
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep portio.IFileReader, path string) (*motion.VmdMotion, error) {
	return commonusecase.LoadMotion(rep, path)
}

// LoadModelWithValidation はモデルを読み込み、結果を返す。
// テクスチャ検証は mu_model_viewer 固有処理のため、ここでは実行しない。
func LoadModelWithValidation(rep portio.IFileReader, path string, validator portio.ITextureValidator) (*ModelLoadResult, error) {
	_ = validator // テクスチャ検証はmu_model_viewer固有のため、このusecaseでは扱わない。
	modelData, err := LoadModel(rep, path)
	if err != nil {
		return nil, err
	}
	return &ModelLoadResult{Model: modelData}, nil
}

// LoadMotionWithMeta はモーションを読み込み、最大フレームを返す。
func LoadMotionWithMeta(rep portio.IFileReader, path string) (*MotionLoadResult, error) {
	return commonusecase.LoadMotionWithMeta(rep, path)
}

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func CanLoadPath(rep portio.IFileReader, path string) bool {
	return commonusecase.CanLoadPath(rep, path)
}
