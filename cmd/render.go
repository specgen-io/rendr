package cmd

import (
	"fmt"
	"github.com/specgen-io/rendr/render"
	"github.com/specgen-io/rendr/values"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

const OutPath = "out"
const Blueprint = "blueprint"
const Set = "set"
const ExtraRoots = "root"
const Values = "values"
const NoInput = "noinput"
const ForceInput = "forceinput"
const NoOverwrites = "nooverwrites"

func init() {
	cobra.OnInitialize()
	cmdNew.Flags().String(Blueprint, "rendr.yaml", `blueprint file inside of the template`)
	cmdNew.Flags().String(OutPath, ".", `path to output rendered template`)
	cmdNew.Flags().StringArray(Set, []string{}, `set arguments overrides in format "arg=value", repeat for setting multiple arguments values`)
	cmdNew.Flags().String(Values, "", `path to arguments values file, could json or yaml`)
	cmdNew.Flags().Bool(NoInput, false, `do not request user input for missing arguments values`)
	cmdNew.Flags().Bool(ForceInput, false, `force user input requests even for noinput arguments`)
	cmdNew.Flags().Bool(NoOverwrites, false, `do not overwrite files with rendered from template`)
	cmdNew.Flags().StringArray(ExtraRoots, []string{}, `extra template root, repeat for setting multiple extra roots`)
}

var cmdNew = &cobra.Command{
	Use:   "rendr <template-url> [flags]",
	Short: "Render template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateUrl := args[0]

		blueprintPath, err := cmd.Flags().GetString(Blueprint)
		failIfError(err, `Failed to get "%s" option`, Blueprint)

		outPath, err := cmd.Flags().GetString(OutPath)
		failIfError(err, `Failed to get "%s" option`, OutPath)

		overrides, err := cmd.Flags().GetStringArray(Set)
		failIfError(err, `Failed to get "%s" option`, Set)

		valuesFilePath, err := cmd.Flags().GetString(Values)
		failIfError(err, `Failed to get "%s" option`, Values)

		noInput, err := cmd.Flags().GetBool(NoInput)
		failIfError(err, `Failed to get "%s" option`, NoInput)

		forceInput, err := cmd.Flags().GetBool(ForceInput)
		failIfError(err, `Failed to get "%s" option`, ForceInput)

		noOverwrites, err := cmd.Flags().GetBool(NoOverwrites)
		failIfError(err, `Failed to get "%s" option`, NoOverwrites)

		extraRoots, err := cmd.Flags().GetStringArray(ExtraRoots)
		failIfError(err, `Failed to get "%s" option`, ExtraRoots)

		inputMode := render.RegularInputMode
		if forceInput {
			inputMode = render.ForceInputMode
		}
		if noInput {
			inputMode = render.NoInputMode
		}

		var valuesData []byte = nil
		if valuesFilePath != "" {
			data, err := ioutil.ReadFile(valuesFilePath)
			failIfError(err, `Failed to read arguments file "%s"`, valuesFilePath)
			valuesData = data
		}
		valuesDataKind := values.JSON
		if strings.HasSuffix(valuesFilePath, ".yaml") || strings.HasSuffix(valuesFilePath, ".yml") {
			valuesDataKind = values.YAML
		}

		templateUrl = normalizeTemplateUrl(templateUrl)
		err = renderTemplate(templateUrl, extraRoots, blueprintPath, outPath, inputMode, &values.ValuesData{valuesDataKind, valuesData}, overrides, !noOverwrites)
		failIfError(err, "Failed to render template")
	},
}

func normalizeTemplateUrl(templateUrl string) string {
	if strings.HasPrefix(templateUrl, "github.com") {
		parts := strings.Split(templateUrl, "/")
		githubSlug := parts[0:3]
		templateUrl = fmt.Sprintf(`https://%s.git`, strings.Join(githubSlug, "/"))
		if len(parts) > 3 {
			pathParts := parts[3:]
			templateUrl = fmt.Sprintf(`%s/%s`, templateUrl, strings.Join(pathParts, "/"))
		}
	}
	return templateUrl
}

func renderTemplate(sourceUrl string, extraRoots []string, blueprintPath string, outPath string, inputMode render.InputMode, valuesData *values.ValuesData, overrides []string, overwriteFiles bool) error {
	template := render.Template{sourceUrl, blueprintPath, extraRoots}
	renderedFiles, err := template.Render(inputMode, valuesData, overrides)
	if err != nil {
		return err
	}

	err = renderedFiles.WriteAll(outPath, overwriteFiles)
	if err != nil {
		return err
	}
	return nil
}

func failIfError(err error, format string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, format, args...)
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func Execute() {
	println(`Running rendr`)
	if err := cmdNew.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
