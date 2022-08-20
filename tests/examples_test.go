package tests

import (
	"fmt"
	"github.com/specgen-io/rendr/render"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var exampleTestCases = []string{
	"simple",
}

func Test_Examples(t *testing.T) {
	examplesPath, err := filepath.Abs("./examples")
	if err != nil {
		t.Fatalf(`failed to get absolute path for "%s": %s`, "./examples", err.Error())
	}

	expectedPath, err := filepath.Abs("./expected")
	if err != nil {
		t.Fatalf(`failed to get absolute path for "%s": %s`, "./expected", err.Error())
	}

	actualPath, err := filepath.Abs("./actual")
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
		t.Logf(`Running test case: %s`, testcase)

		templatePath := filepath.Join(examplesPath, testcase)

		var valuesJsonData []byte = nil
		valuesJsonPath := filepath.Join(expectedPath, fmt.Sprintf(`%s.json`, testcase))
		if render.Exists(valuesJsonPath) {
			data, err := ioutil.ReadFile(valuesJsonPath)
			if err != nil {
				t.Fatalf(`failed to read file "%s": %s`, valuesJsonPath, err.Error())
			}
			valuesJsonData = data
		}

		var overrides []string = nil
		//overridesPath := filepath.Join(expectedPath, fmt.Sprintf(`%s.overrides`, testcase))
		//if render.Exists(overridesPath) {
		//
		//}

		outPath, err := os.MkdirTemp(actualPath, testcase)
		if err != nil {
			t.Fatalf(`failed to get temp folder path: %s`, err.Error())
		}
		err = os.Chmod(outPath, 0755)
		if err != nil {
			t.Fatalf(`failed to change mode for folder "%s": %s`, outPath, err.Error())
		}
		err = RenderExampleTemplate(templatePath, valuesJsonData, overrides, outPath)
		if err != nil {
			t.Fatalf(`failed to render template: %s`, err.Error())
		}

		expectedCasePath := filepath.Join(expectedPath, testcase)

		assert.Assert(t, fs.Equal(outPath, fs.ManifestFromDir(t, expectedCasePath)))
	}
}

func RenderExampleTemplate(templatePath string, valuesJson []byte, overrides []string, outPath string) error {
	templateUrl := fmt.Sprintf(`file:///%s`, templatePath)
	template := render.Template{templateUrl, "", "rendr.yaml"}
	renderedFiles, err := template.Render(render.NoInputMode, valuesJson, overrides)
	if err != nil {
		return err
	}
	err = renderedFiles.WriteAll(outPath, true)
	if err != nil {
		return err
	}
	return nil
}
