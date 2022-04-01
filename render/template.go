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
	"github.com/specgen-io/rendr/files"
	"io/fs"
	"path"
	"strings"
)

type Template struct {
	RepoUrl       string
	Path          string
	BlueprintPath string
}

func (t Template) Render(outPath string, noInput bool, forceInput bool, valuesJsonData []byte, overridesKeysValues []string) ([]files.Text, error) {
	filesystem, err := getFilesystem(t.RepoUrl)
	if err != nil {
		return nil, err
	}

	blueprint, err := t.LoadBlueprint(filesystem)
	if err != nil {
		return nil, err
	}

	argsValues, err := t.GetArgsValues(blueprint.Args, noInput, forceInput, valuesJsonData, overridesKeysValues)
	if err != nil {
		return nil, err
	}

	result := []files.Text{}

	for _, root := range blueprint.Roots {
		rootFiles, err := t.RenderRoot(filesystem, root, blueprint, argsValues, outPath)
		if err != nil {
			return nil, err
		}
		result = append(result, rootFiles...)
	}

	return result, err
}

func (t Template) RenderRoot(
	filesystem billy.Filesystem,
	root string,
	blueprint *blueprint.Blueprint,
	argsValues blueprint.ArgsValues,
	outPath string) ([]files.Text, error) {

	result := []files.Text{}

	staticFiles, err := t.getFiles(filesystem, root, blueprint.StaticPaths, blueprint.IgnorePaths)
	if err != nil {
		return nil, err
	}

	templateFiles, err := t.getFiles(filesystem, root, nil, blueprint.IgnorePaths)
	if err != nil {
		return nil, err
	}

	staticResults, err := renderStaticFiles(staticFiles, outPath)
	if err != nil {
		return nil, err
	}
	result = append(result, staticResults...)

	renderedResults, err := renderFiles(templateFiles, outPath, argsValues)
	if err != nil {
		return nil, err
	}
	result = append(result, renderedResults...)

	return result, nil
}

func renderFiles(templateFiles []files.Text, outPath string, argsValues blueprint.ArgsValues) ([]files.Text, error) {
	result := []files.Text{}
	for _, templateFile := range templateFiles {
		renderedFile, err := renderFile(&templateFile, outPath, argsValues)
		if err != nil {
			return nil, fmt.Errorf(`template "%s" returned error: %s`, templateFile.Path, err.Error())
		}
		if renderedFile != nil {
			result = append(result, *renderedFile)
		}
	}
	return result, nil
}

func renderFile(templateFile *files.Text, outPath string, argsValues blueprint.ArgsValues) (*files.Text, error) {
	templatePath := templateFile.Path

	renderedPath, err := renderPath(templatePath, argsValues)
	if err != nil {
		return nil, err
	}

	if renderedPath == nil {
		return nil, nil
	}

	content, err := render(templateFile.Content, argsValues)
	if err != nil {
		return nil, err
	}

	return &files.Text{path.Join(outPath, *renderedPath), content}, nil
}

func renderStaticFiles(staticFiles []files.Text, outPath string) ([]files.Text, error) {
	result := []files.Text{}
	for _, theFile := range staticFiles {
		content := theFile.Content
		result = append(result, files.Text{path.Join(outPath, theFile.Path), content})
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
	if result.StaticPaths == nil {
		result.StaticPaths = blueprint.PathPrefixArray{}
	}
	return result, nil
}

func (t Template) getFiles(filesystem billy.Filesystem, rootPath string, includeOnlyPrefixes blueprint.PathPrefixArray, excludePrefixes blueprint.PathPrefixArray) ([]files.Text, error) {
	result := []files.Text{}
	rootFullPath := path.Join(t.Path, rootPath)
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
			file := files.Text{filepath, string(data)}
			if includeOnlyPrefixes == nil || includeOnlyPrefixes.Matches(filepath) {
				result = append(result, file)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
