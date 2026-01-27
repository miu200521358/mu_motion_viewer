package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	portio "github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

// ModelLoadResult はモデル読み込み結果を表す。
type ModelLoadResult = commonusecase.ModelLoadResult

// MotionLoadResult はモーション読み込み結果を表す。
type MotionLoadResult = commonusecase.MotionLoadResult

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep portio.IFileReader, path string) (*model.PmxModel, error) {
	return commonusecase.LoadModel(rep, path)
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep portio.IFileReader, path string) (*motion.VmdMotion, error) {
	return commonusecase.LoadMotion(rep, path)
}

// LoadModelWithValidation はモデルを読み込み、必要に応じてテクスチャ検証を行う。
func LoadModelWithValidation(rep portio.IFileReader, path string, validator portio.ITextureValidator) (*ModelLoadResult, error) {
	return commonusecase.LoadModelWithValidation(rep, path, validator)
}

// LoadMotionWithMeta はモーションを読み込み、最大フレームを返す。
func LoadMotionWithMeta(rep portio.IFileReader, path string) (*MotionLoadResult, error) {
	return commonusecase.LoadMotionWithMeta(rep, path)
}

// CanLoadPath はリポジトリが指定パスを読み込み可能か判定する。
func CanLoadPath(rep portio.IFileReader, path string) bool {
	return commonusecase.CanLoadPath(rep, path)
}
