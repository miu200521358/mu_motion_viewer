// 指示: miu200521358
package minteractor

import "github.com/miu200521358/mu_tree_viewer/pkg/usecase/port/moutput"

// TreeViewerUsecaseDeps はツリービューア用ユースケースの依存を表す。
type TreeViewerUsecaseDeps struct {
	ModelReader  moutput.IFileReader
	MotionReader moutput.IFileReader
}

// TreeViewerUsecase はツリービューアの入出力処理をまとめたユースケースを表す。
type TreeViewerUsecase struct {
	modelReader  moutput.IFileReader
	motionReader moutput.IFileReader
}

// NewTreeViewerUsecase はツリービューア用ユースケースを生成する。
func NewTreeViewerUsecase(deps TreeViewerUsecaseDeps) *TreeViewerUsecase {
	return &TreeViewerUsecase{
		modelReader:  deps.ModelReader,
		motionReader: deps.MotionReader,
	}
}
