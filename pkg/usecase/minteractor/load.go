// 指示: miu200521358
package minteractor

import (
	commonusecase "github.com/miu200521358/mlib_go/pkg/usecase"
	"github.com/miu200521358/mu_tree_viewer/pkg/usecase/port/moutput"
)

// LoadModel はモデルを読み込み、結果を返す。
func (uc *TreeViewerUsecase) LoadModel(rep moutput.IFileReader, path string) (*ModelLoadResult, error) {
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
func (uc *TreeViewerUsecase) LoadMotion(rep moutput.IFileReader, path string) (*MotionLoadResult, error) {
	repo := rep
	if repo == nil {
		repo = uc.motionReader
	}
	result, err := commonusecase.LoadMotionWithMeta(repo, path)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return &MotionLoadResult{Motion: result.Motion, MaxFrame: result.MaxFrame}, nil
}

// CanLoadModelPath はモデルの読み込み可否を判定する。
func (uc *TreeViewerUsecase) CanLoadModelPath(path string) bool {
	return commonusecase.CanLoadPath(uc.modelReader, path)
}
