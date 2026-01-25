package usecase

import (
	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"

	"github.com/miu200521358/mu_motion_viewer/pkg/workflow"
)

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep io_common.IFileReader, path string) (*model.PmxModel, error) {
	return workflow.LoadModel(rep, path)
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep io_common.IFileReader, path string) (*motion.VmdMotion, error) {
	return workflow.LoadMotion(rep, path)
}
