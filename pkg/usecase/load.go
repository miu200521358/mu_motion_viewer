package usecase

import (
	"fmt"
	"os"
	"strings"

	"github.com/miu200521358/mlib_go/pkg/config/mconfig"
)

func LoadModelMotionPath() (modelPath, motionPath string, err error) {
	if len(os.Args) <= 1 {
		return "", "", fmt.Errorf("no motion file")
	}

	modelPaths := mconfig.LoadUserConfig("pmx")
	for _, path := range modelPaths {
		if strings.LastIndex(strings.ToLower(path), ".pmx") == len(path)-4 {
			modelPath = path
			break
		}
	}

	if modelPath == "" {
		return "", "", fmt.Errorf("no model file")
	}

	motionPath = os.Args[1]
	return modelPath, motionPath, nil
}
