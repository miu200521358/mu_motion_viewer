package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
)

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep io_common.IFileReader, path string) (*model.PmxModel, error) {
	return commonusecase.LoadModel(rep, path)
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep io_common.IFileReader, path string) (*motion.VmdMotion, error) {
	return commonusecase.LoadMotion(rep, path)
}
