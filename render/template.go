package render

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/specgen-io/rendr/blueprint"
	"github.com/specgen-io/rendr/values"
	"io/fs"
	"path"
	"strings"
)

type Template struct {
	Source        string
	Path          string
	BlueprintPath string
}

type InputMode string

const (
	RegularInputMode InputMode = "regular"
	NoInputMode      InputMode = "no"
	ForceInputMode   InputMode = "force"
)

func (t Template) Render(inputMode InputMode, valuesJsonData []byte, overridesKeysValues []string) (Files, error) {
	filesystem, err := getFilesystem(t.Source)
	if err != nil {
		return nil, err
	}

	blueprint, err := t.LoadBlueprint(filesystem)
	if err != nil {
		return nil, err
	}

	argsValues, err := t.GetArgsValues(blueprint.Args, inputMode, valuesJsonData, overridesKeysValues)
	if err != nil {
		return nil, err
	}

	files := []File{}

	for _, root := range blueprint.Roots {
		rootFiles, err := t.RenderRoot(filesystem, root, blueprint, argsValues)
		if err != nil {
			return nil, err
		}
		files = append(files, rootFiles...)
	}

	for source, target := range blueprint.Rename {
		for i := range files {
			path := files[i].Path
			if strings.HasPrefix(path, source) {
				files[i].Path = strings.Replace(path, source, target, 1)
			}
		}
	}

	return files, err
}

func (t Template) RenderRoot(
	filesystem billy.Filesystem,
	root string,
	blueprint *blueprint.Blueprint,
	argsValues values.ArgsValues) ([]File, error) {

	renderFullPath := path.Join(t.Path, root)
	templateFiles, err := getFiles(filesystem, renderFullPath, blueprint.IgnorePaths, blueprint.ExecutablePaths, blueprint.StaticPaths)
	if err != nil {
		return nil, err
	}

	renderedFiles, err := renderFiles(templateFiles, argsValues)
	if err != nil {
		return nil, err
	}

	return renderedFiles, nil
}

func (t Template) LoadBlueprint(filesystem billy.Filesystem) (*blueprint.Blueprint, error) {
	blueprintFullpath := path.Join(t.Path, t.BlueprintPath)
	data, err := util.ReadFile(filesystem, blueprintFullpath)
	if err != nil {
		return nil, err
	}
	result, err := blueprint.Read(string(data))
	if err != nil {
		return nil, err
	}
	result.IgnorePaths = append(result.IgnorePaths, t.BlueprintPath)
	if result == nil || len(result.Roots) == 0 {
		result.Roots = []string{"."}
	}
	if result.Rename == nil {
		result.Rename = map[string]string{}
	}
	return result, nil
}

func getFilesystem(repoUrl string) (billy.Filesystem, error) {
	if strings.HasPrefix(repoUrl, "file:///") {
		repoPath := strings.TrimPrefix(repoUrl, "file:///")
		return osfs.New(repoPath), nil
	} else {
		filesystem := memfs.New()
		_, err := git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{URL: repoUrl})
		if err != nil {
			return nil, err
		}
		return filesystem, nil
	}
}

func getFiles(filesystem billy.Filesystem, rootFullPath string, excludePrefixes blueprint.PathArray, executablePaths blueprint.PathArray, staticPaths blueprint.PathArray) ([]File, error) {
	result := []File{}
	err := Walk(filesystem, rootFullPath, func(itempath string, info fs.FileInfo, err error) error {
		filepath := strings.TrimPrefix(strings.TrimPrefix(itempath, rootFullPath), "/")
		if excludePrefixes.Matches(filepath) {
			return nil
		}
		if !info.IsDir() {
			data, err := util.ReadFile(filesystem, itempath)
			if err != nil {
				return nil
			}
			executable := executablePaths.Contains(filepath)
			static := staticPaths.Matches(filepath)
			template := !executable && !static
			file := File{filepath, string(data), executable, template}
			result = append(result, file)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func renderFile(sourceFile *File, argsValues values.ArgsValues) (*File, error) {
	renderedPath, err := renderPath(sourceFile.Path, argsValues)
	if err != nil {
		return nil, err
	}
	if renderedPath == nil {
		return nil, nil
	}

	content := sourceFile.Content
	if sourceFile.Template {
		content, err = values.Render(content, argsValues)
		if err != nil {
			return nil, err
		}
	}

	return &File{*renderedPath, content, sourceFile.Executable, false}, nil
}

func renderFiles(templateFiles []File, argsValues values.ArgsValues) ([]File, error) {
	result := []File{}
	for _, templateFile := range templateFiles {
		renderedFile, err := renderFile(&templateFile, argsValues)
		if err != nil {
			return nil, fmt.Errorf(`template "%s" returned error: %s`, templateFile.Path, err.Error())
		}
		if renderedFile != nil {
			result = append(result, *renderedFile)
		}
	}
	return result, nil
}

func renderPath(templatePath string, argsValues values.ArgsValues) (*string, error) {
	parts := strings.Split(templatePath, "/")
	resultParts := []string{}
	for _, part := range parts {
		resultPart, err := values.RenderShort(part, argsValues)
		if err != nil {
			return nil, err
		}
		if resultPart == nil {
			return nil, nil
		}
		if *resultPart != "" {
			resultParts = append(resultParts, *resultPart)
		}
	}
	result := strings.Join(resultParts, "/")
	return &result, nil
}
