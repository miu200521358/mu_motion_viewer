package workflow

import (
	"fmt"

	"github.com/miu200521358/mlib_go/pkg/adapter/io_common"
	"github.com/miu200521358/mlib_go/pkg/domain/model"
	"github.com/miu200521358/mlib_go/pkg/domain/motion"
)

// LoadModel はモデルを読み込み、型を検証して返す。
func LoadModel(rep io_common.IFileReader, path string) (*model.PmxModel, error) {
	if path == "" {
		return nil, nil
	}
	if rep == nil {
		return nil, fmt.Errorf("モデル読み込みリポジトリがありません")
	}
	data, err := rep.Load(path)
	if err != nil {
		return nil, err
	}
	modelData, ok := data.(*model.PmxModel)
	if !ok {
		return nil, fmt.Errorf("モデル形式が不正です")
	}
	return modelData, nil
}

// LoadMotion はモーションを読み込み、型を検証して返す。
func LoadMotion(rep io_common.IFileReader, path string) (*motion.VmdMotion, error) {
	if path == "" {
		return nil, nil
	}
	if rep == nil {
		return nil, fmt.Errorf("モーション読み込みリポジトリがありません")
	}
	data, err := rep.Load(path)
	if err != nil {
		return nil, err
	}
	motionData, ok := data.(*motion.VmdMotion)
	if !ok {
		return nil, fmt.Errorf("モーション形式が不正です")
	}
	return motionData, nil
}
