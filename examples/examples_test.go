package examples

import (
	"bufio"
	"fmt"
	"github.com/specgen-io/rendr/render"
	"github.com/specgen-io/rendr/values"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/fs"
	"os"
	"path/filepath"
	"testing"
)

var exampleTestCases = []ExampleTestCase{
	{"simple", "json_values"},
	{"simple", "yaml_values"},
	{"simple", "override_values"},
	{"folders", "folders"},
}

type ExampleTestCase struct {
	Template string
	Expected string
}

func Test_Examples(t *testing.T) {
	examplesPath, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf(`failed to get absolute path for "%s": %s`, "./examples", err.Error())
	}

	expectedPath, err := filepath.Abs("./_expected")
	if err != nil {
		t.Fatalf(`failed to get absolute path for "%s": %s`, "./expected", err.Error())
	}

	actualPath, err := filepath.Abs("./_actual")
	if err != nil {
		t.Fatalf(`failed to get absolute path for "%s": %s`, "./actual", err.Error())
	}
	if !render.Exists(actualPath) {
		err := os.MkdirAll(actualPath, os.ModePerm)
		if err != nil {
			t.Fatalf(`failed to create folder "%s": %s`, actualPath, err.Error())
		}
	}

	for _, testcase := range exampleTestCases {
		t.Logf(`Running test case: %s`, testcase.Expected)

		templatePath := filepath.Join(examplesPath, testcase.Template)

		expectedCasePath := filepath.Join(expectedPath, testcase.Expected)

		valuesPath := filepath.Join(expectedCasePath, `values.yaml`)
		if !render.Exists(valuesPath) {
			valuesPath = filepath.Join(expectedCasePath, `values.json`)
		}
		var valuesData *values.ValuesData = nil
		if render.Exists(valuesPath) {
			valuesData, err = values.LoadValuesFile(valuesPath)
			if err != nil {
				t.Fatalf(`failed to load values data "%s": %s`, valuesPath, err.Error())
			}
		}

		var overrides []string = nil
		overridesPath := filepath.Join(expectedCasePath, `values.overrides`)
		if render.Exists(overridesPath) {
			overrides = []string{}
			file, err := os.Open(overridesPath)
			if err != nil {
				t.Fatalf(`failed to read file "%s": %s`, overridesPath, err.Error())
			}

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				overrides = append(overrides, scanner.Text())
			}
			file.Close()
		}

		outPath, err := os.MkdirTemp(actualPath, testcase.Expected)
		if err != nil {
			t.Fatalf(`failed to get temp folder path: %s`, err.Error())
		}
		err = os.Chmod(outPath, 0755)
		if err != nil {
			t.Fatalf(`failed to change mode for folder "%s": %s`, outPath, err.Error())
		}
		err = RenderExampleTemplate(templatePath, valuesData, overrides, outPath)
		if err != nil {
			t.Fatalf(`failed to render template: %s`, err.Error())
		}

		expectedFilesPath := filepath.Join(expectedCasePath, `files`)

		assert.Assert(t, fs.Equal(outPath, fs.ManifestFromDir(t, expectedFilesPath)))
	}
}

func RenderExampleTemplate(templatePath string, valuesData *values.ValuesData, overrides []string, outPath string) error {
	templateUrl := fmt.Sprintf(`file:///%s`, templatePath)
	template := render.Template{templateUrl, "rendr.yaml", nil}
	renderedFiles, err := template.Render(render.NoInputMode, valuesData, overrides)
	if err != nil {
		return err
	}
	err = renderedFiles.WriteAll(outPath, true)
	if err != nil {
		return err
	}
	return nil
}
