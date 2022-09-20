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
	BlueprintPath string
	ExtraRoots    []string
}

type InputMode string

const (
	RegularInputMode InputMode = "regular"
	NoInputMode      InputMode = "no"
	ForceInputMode   InputMode = "force"
)

func (t *Template) Render(inputMode InputMode, valuesJsonData []byte, overridesKeysValues []string) (Files, error) {
	blueprint, err := t.LoadBlueprint()
	if err != nil {
		return nil, err
	}

	argsValues, err := t.GetArgsValues(blueprint.Args, inputMode, valuesJsonData, overridesKeysValues)
	if err != nil {
		return nil, err
	}

	files := []File{}

	roots := t.GetRoots(blueprint)
	for _, root := range roots {
		rootFiles, err := renderRoot(root, blueprint, argsValues)
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

type Root struct {
	Source string
	Path   string
}

func (t *Template) GetRoots(blueprint *blueprint.Blueprint) []string {
	result := []string{}
	for _, rootPath := range blueprint.Roots {
		rootFullPath := t.Source
		if rootPath != "." {
			rootFullPath = fmt.Sprintf("%s/%s", rootFullPath, rootPath)
		}
		result = append(result, rootFullPath)
	}
	if t.ExtraRoots != nil {
		result = append(result, t.ExtraRoots...)
	}
	return result
}

func renderRoot(
	rootUrl string,
	blueprint *blueprint.Blueprint,
	argsValues values.ArgsValues) ([]File, error) {

	source, rootPath := splitSource(rootUrl)
	filesystem, err := getFilesystem(source)
	if err != nil {
		return nil, err
	}

	templateFiles, err := getFiles(filesystem, rootPath, blueprint.IgnorePaths, blueprint.ExecutablePaths, blueprint.StaticPaths)
	if err != nil {
		return nil, err
	}

	renderedFiles, err := renderFiles(templateFiles, argsValues)
	if err != nil {
		return nil, err
	}

	return renderedFiles, nil
}

func splitSource(sourceUrl string) (string, string) {
	if strings.HasPrefix(sourceUrl, "file:///") || strings.HasSuffix(sourceUrl, ".git") {
		return sourceUrl, ""
	}
	parts := strings.Split(sourceUrl, ".git/")
	path := parts[1]
	source := sourceUrl[:len(sourceUrl)-1-len(path)]
	return source, path
}

func (t *Template) LoadBlueprint() (*blueprint.Blueprint, error) {
	source, sourcePath := splitSource(t.Source)
	filesystem, err := getFilesystem(source)
	if err != nil {
		return nil, err
	}
	blueprintFullpath := path.Join(sourcePath, t.BlueprintPath)
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

var filesystems = make(map[string]billy.Filesystem)

func getFilesystem(url string) (billy.Filesystem, error) {
	if filesystem, found := filesystems[url]; found {
		return filesystem, nil
	}
	if strings.HasPrefix(url, "file:///") {
		repoPath := strings.TrimPrefix(url, "file:///")
		filesystem := osfs.New(repoPath)
		filesystems[url] = filesystem
		return filesystem, nil
	} else {
		filesystem := memfs.New()
		_, err := git.Clone(memory.NewStorage(), filesystem, &git.CloneOptions{URL: url})
		if err != nil {
			return nil, err
		}
		filesystems[url] = filesystem
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
