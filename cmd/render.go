package cmd

import (
	"fmt"
	"github.com/specgen-io/rendr/console"
	"github.com/specgen-io/rendr/render"
	"github.com/specgen-io/rendr/values"
	"github.com/spf13/cobra"
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
const Verbose = "verbose"

func init() {
	cobra.OnInitialize()
	cmdRoot.Flags().String(Blueprint, "rendr.yaml", `blueprint file inside of the template`)
	cmdRoot.Flags().String(OutPath, ".", `path to output rendered template`)
	cmdRoot.Flags().StringArray(Set, []string{}, `set arguments overrides in format "arg=value", repeat for setting multiple arguments values`)
	cmdRoot.Flags().String(Values, "", `path to arguments values file, could json or yaml`)
	cmdRoot.Flags().Bool(NoInput, false, `do not request user input for missing arguments values`)
	cmdRoot.Flags().Bool(ForceInput, false, `force user input requests even for noinput arguments`)
	cmdRoot.Flags().Bool(NoOverwrites, false, `do not overwrite files with rendered from template`)
	cmdRoot.Flags().StringArray(ExtraRoots, []string{}, `extra template root, repeat for setting multiple extra roots`)
	cmdRoot.Flags().Bool(Verbose, false, `print more logging`)
}

var cmdRoot = &cobra.Command{
	Use:   "rendr <template-url> [flags]",
	Short: "Render template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, err := cmd.Flags().GetBool(Verbose)
		failIfError(err, `Failed to get "%s" option`, Verbose)
		if verbose {
			console.Level = console.VerboseLevel
		}

		console.Verbose("Running rendr")

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

		valuesData, err := values.LoadValuesFile(valuesFilePath)
		failIfError(err, `Failed to load values file "%s"`, valuesFilePath)

		templateUrl = normalizeTemplateUrl(templateUrl)
		err = renderTemplate(templateUrl, extraRoots, blueprintPath, outPath, inputMode, valuesData, overrides, !noOverwrites)
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
	console.Error(err, format, args...)
	if err != nil {
		os.Exit(1)
	}
}

func Execute() {
	err := cmdRoot.Execute()
	failIfError(err, "Failed to run rendr tool")
}
