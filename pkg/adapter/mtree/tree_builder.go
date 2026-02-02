// 指示: miu200521358
package mtree

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	rootLabelFormat = "【%s】"
)

var defaultModelExtensions = []string{".pmx", ".pmd", ".x"}

// Node はモデルツリーのノードを表す。
type Node struct {
	Name     string
	Path     string
	IsFile   bool
	Children []*Node
}

// BuildModelTree はモデルファイル用のツリーを構築する。
func BuildModelTree(paths []string) ([]*Node, error) {
	return BuildTree(paths, defaultModelExtensions)
}

// BuildTree は指定拡張子に一致するファイルのみでツリーを構築する。
func BuildTree(paths []string, extensions []string) ([]*Node, error) {
	cleaned := normalizePaths(paths)
	if len(cleaned) == 0 {
		return []*Node{}, nil
	}

	var roots []*Node
	var errs []error
	for _, path := range cleaned {
		absPath, err := filepath.Abs(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("パスの絶対化に失敗しました: %w", err))
			continue
		}
		info, err := os.Stat(absPath)
		if err != nil {
			errs = append(errs, fmt.Errorf("フォルダ情報の取得に失敗しました: %w", err))
			continue
		}
		if !info.IsDir() {
			continue
		}
		node, ok, buildErr := buildFolderNode(absPath, extensions)
		if buildErr != nil {
			errs = append(errs, buildErr)
			continue
		}
		if ok {
			node.Name = fmt.Sprintf(rootLabelFormat, absPath)
			roots = append(roots, node)
		}
	}

	sort.SliceStable(roots, func(i, j int) bool {
		return roots[i].Path < roots[j].Path
	})

	return roots, errors.Join(errs...)
}

// buildFolderNode はフォルダ配下のツリーを構築する。
func buildFolderNode(path string, extensions []string) (*Node, bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, false, fmt.Errorf("フォルダ一覧の取得に失敗しました: %w", err)
	}

	node := &Node{
		Name:   filepath.Base(path),
		Path:   path,
		IsFile: false,
	}

	var children []*Node
	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(path, name)
		if entry.IsDir() {
			child, ok, childErr := buildFolderNode(fullPath, extensions)
			if childErr != nil {
				return nil, false, childErr
			}
			if ok {
				children = append(children, child)
			}
			continue
		}
		if !hasExtension(name, extensions) {
			continue
		}
		children = append(children, &Node{
			Name:   name,
			Path:   fullPath,
			IsFile: true,
		})
	}

	sort.SliceStable(children, func(i, j int) bool {
		left := strings.ToLower(children[i].Name)
		right := strings.ToLower(children[j].Name)
		if left == right {
			return children[i].Name < children[j].Name
		}
		return left < right
	})

	node.Children = children
	return node, len(children) > 0, nil
}

// normalizePaths は入力パスを正規化して重複排除する。
func normalizePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		cleaned := cleanPath(p)
		if cleaned == "" {
			continue
		}
		if _, ok := seen[cleaned]; ok {
			continue
		}
		seen[cleaned] = struct{}{}
		result = append(result, cleaned)
	}
	return result
}

// cleanPath は入力パスを正規化して返す。
func cleanPath(path string) string {
	trimmed := strings.TrimSpace(path)
	trimmed = strings.Trim(trimmed, "\"")
	trimmed = strings.Trim(trimmed, "'")
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}

// hasExtension は拡張子判定を大文字小文字を区別せずに行う。
func hasExtension(name string, extensions []string) bool {
	if name == "" {
		return false
	}
	ext := filepath.Ext(name)
	if ext == "" {
		return false
	}
	for _, allowed := range extensions {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}
