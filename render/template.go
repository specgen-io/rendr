package render

import (
	"github.com/cbroglie/mustache"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/specgen-io/rendr/blueprint"
	"github.com/specgen-io/rendr/files"
	"github.com/specgen-io/rendr/input"
	"io/fs"
	"path"
	"strings"
)

type Template struct {
	RepoUrl       string
	Path          string
	BlueprintPath string
}

func (t Template) Render(outPath string, allowInput bool, valuesJsonData []byte, overridesKeysValues []string) ([]files.Text, error) {
	filesystem, err := getFilesystem(t.RepoUrl)
	if err != nil {
		return nil, err
	}

	theBlueprint, err := t.loadBlueprint(filesystem)
	if err != nil {
		return nil, err
	}

	staticFiles, templateFiles, err := t.load(filesystem, theBlueprint)
	if err != nil {
		return nil, err
	}

	argsValues := blueprint.ArgsValues{}

	if valuesJsonData != nil {
		argsValues, err = blueprint.ReadValuesJson(theBlueprint.Args, valuesJsonData)
		if err != nil {
			return nil, err
		}
	}

	if overridesKeysValues != nil {
		overridesValues, err := blueprint.ParseValues(theBlueprint.Args, overridesKeysValues)
		if err != nil {
			return nil, err
		}
		argsValues, err = blueprint.OverrideValues(theBlueprint.Args, argsValues, overridesValues)
		if err != nil {
			return nil, err
		}
	}

	argsInput := input.NoInput
	if allowInput {
		argsInput = input.Survey
	}
	argsValues, err = blueprint.GetValues(theBlueprint.Args, true, argsValues, argsInput)
	if err != nil {
		return nil, err
	}

	argsValues = blueprint.EnrichValues(theBlueprint.Args, argsValues)

	staticResults, err := renderFiles(staticFiles, outPath, argsValues, true)
	if err != nil {
		return nil, err
	}

	renderedResults, err := renderFiles(templateFiles, outPath, argsValues, false)
	if err != nil {
		return nil, err
	}

	return append(staticResults, renderedResults...), err
}

func renderFiles(templateFiles []files.Text, outPath string, argsValues blueprint.ArgsValues, isStaticFile bool) ([]files.Text, error) {
	result := []files.Text{}
	for _, templateFile := range templateFiles {
		content := templateFile.Content
		if !isStaticFile {
			mustache.AllowMissingVariables = false
			renderedContent, err := mustache.Render(content, argsValues)
			if err != nil {
				return nil, err
			}
			content = renderedContent
		}
		result = append(result, files.Text{path.Join(outPath, templateFile.Path), content})
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

func (t Template) loadBlueprint(filesystem billy.Filesystem) (*blueprint.Blueprint, error) {
	blueprintFullpath := path.Join(t.Path, t.BlueprintPath)
	data, err := util.ReadFile(filesystem, blueprintFullpath)
	if err != nil {
		return nil, err
	}
	theBlueprint, err := blueprint.Read(string(data))
	if err != nil {
		return nil, err
	}
	theBlueprint.IgnorePaths = append(theBlueprint.IgnorePaths, t.BlueprintPath)
	return theBlueprint, nil
}

func (t Template) load(filesystem billy.Filesystem, blueprint *blueprint.Blueprint) ([]files.Text, []files.Text, error) {
	templateFiles := []files.Text{}
	staticFiles := []files.Text{}
	err := Walk(filesystem, t.Path, func(itempath string, info fs.FileInfo, err error) error {
		filepath := strings.TrimPrefix(strings.TrimPrefix(itempath, t.Path), "/")
		if blueprint.IgnorePaths.Matches(filepath) {
			return nil
		}
		if !info.IsDir() {
			data, err := util.ReadFile(filesystem, itempath)
			if err != nil {
				return nil
			}
			file := files.Text{filepath, string(data)}
			if blueprint.StaticPaths.Matches(filepath) {
				staticFiles = append(staticFiles, file)
			} else {
				templateFiles = append(templateFiles, file)
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return staticFiles, templateFiles, nil
}
