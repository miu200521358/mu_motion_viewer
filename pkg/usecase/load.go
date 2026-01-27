package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	portio "github.com/miu200521358/mlib_go/pkg/usecase/port/io"
)

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep portio.IFileReader, path string) (*model.PmxModel, error) {
	return commonusecase.LoadModel(rep, path)
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep portio.IFileReader, path string) (*motion.VmdMotion, error) {
	return commonusecase.LoadMotion(rep, path)
}
