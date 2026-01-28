// 指示: miu200521358
package minteractor

import "github.com/miu200521358/mu_motion_viewer/pkg/usecase/port/moutput"

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
