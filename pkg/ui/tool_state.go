package ui

import (
	"errors"
	"os"
	"sort"

	"github.com/miu200521358/mlib_go/pkg/domain/pmx"
	"github.com/miu200521358/mlib_go/pkg/domain/vmd"
	"github.com/miu200521358/mlib_go/pkg/infrastructure/repository"
	"github.com/miu200521358/mlib_go/pkg/interface/app"
	"github.com/miu200521358/mlib_go/pkg/interface/controller"
	"github.com/miu200521358/mlib_go/pkg/interface/controller/widget"
	"github.com/miu200521358/mlib_go/pkg/mutils/mconfig"
)

type ToolState struct {
	App                    *app.MApp
	ControlWindow          *controller.ControlWindow
	ConfigTab              *widget.MTabPage
	PmxPicker              *widget.FilePicker
	Player                 *widget.MotionPlayer
	ModelPath              string
	Model                  *pmx.PmxModel
	Motion                 *vmd.VmdMotion
	ActiveExistBoneNames   []string // モーションに存在して、かつモデルにも存在するボーン名
	ActiveMissingBoneNames []string // モーションに存在して、かつモデルには存在しないボーン名
}

func NewToolState(app *app.MApp, controlWindow *controller.ControlWindow) *ToolState {

	toolState := &ToolState{
		App:           app,
		ControlWindow: controlWindow,
		Motion:        vmd.NewVmdMotion(""),
	}

	modelPaths := mconfig.LoadUserConfig("pmx")
	if len(modelPaths) > 0 {
		modelPath := modelPaths[0]
		rep := repository.NewPmxRepository()
		if isOk, err := rep.CanLoad(modelPath); isOk && err == nil {
			if data, err := rep.Load(modelPath); err == nil {
				toolState.ModelPath = modelPath
				toolState.Model = data.(*pmx.PmxModel)
			} else {
				widget.RaiseError(err)
			}
		} else if err != nil {
			widget.RaiseError(err)
		} else {
			widget.RaiseError(errors.New("unknown error"))
		}
	}

	if toolState.Model != nil && len(os.Args) > 1 {
		motionPath := os.Args[1]
		rep := repository.NewVmdVpdRepository()
		if isOk, err := rep.CanLoad(motionPath); isOk && err == nil {
			if data, err := rep.Load(motionPath); err == nil {
				toolState.Motion = data.(*vmd.VmdMotion)

				for _, boneName := range toolState.Motion.BoneFrames.Names() {
					if toolState.Motion.BoneFrames.ContainsActive(boneName) {
						if toolState.Model.Bones.ContainsByName(boneName) {
							toolState.ActiveExistBoneNames = append(toolState.ActiveExistBoneNames, boneName)
						} else {
							toolState.ActiveMissingBoneNames = append(toolState.ActiveMissingBoneNames, boneName)
						}
					}
				}

				sort.Strings(toolState.ActiveExistBoneNames)
				sort.Strings(toolState.ActiveMissingBoneNames)

				toolState.App.SetFuncGetModels(func() [][]*pmx.PmxModel {
					return [][]*pmx.PmxModel{{toolState.Model}}
				})

				toolState.App.SetFuncGetMotions(func() [][]*vmd.VmdMotion {
					return [][]*vmd.VmdMotion{{toolState.Motion}}
				})
			} else {
				widget.RaiseError(err)
			}
		} else if err != nil {
			widget.RaiseError(err)
		} else {
			widget.RaiseError(errors.New("unknown error"))
		}
	}

	newConfigTab(controlWindow, toolState)

	return toolState
}

type loadPmxResult struct {
	model *pmx.PmxModel
	err   error
}

func (toolState *ToolState) onPlay(playing bool) {
	toolState.SetEnabled(!playing)
}

func (toolState *ToolState) SetEnabled(enabled bool) {
	toolState.PmxPicker.SetEnabled(enabled)
}
