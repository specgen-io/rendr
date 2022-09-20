package main

import (
	"flag"
	"fmt"
	"github.com/specgen-io/rendr/render"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var stderr = log.New(os.Stderr, "", 0)
var stdout = log.New(os.Stdout, "", 0)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: rendr <command> [<parameters>] [<options>]\n")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println(`  render`)
		fmt.Println(`        Renders template.`)
		fmt.Println(`  build`)
		fmt.Println(`        Prints build command for rendered project.`)
		fmt.Println(`  help`)
		fmt.Println(`        Prints this help message.`)
	}

	cmdRender := CmdRender()

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "render":
		cmdRender.Execute(args)
	case "help":
		flag.Usage()
	default:
		fmt.Printf("Expected commands: 'render', 'build', 'help', got: '%s'\n", command)
		os.Exit(1)
	}
	os.Exit(0)
}

func failIfError(err error, format string, args ...interface{}) {
	if err != nil {
		message := fmt.Sprintf(format, args...) + fmt.Sprintf(", %s.", err.Error())
		stderr.Println(message)
		os.Exit(1)
	}
}

type stringArray []string

func (o *stringArray) String() string {
	return strings.Join(*o, ", ")
}

func (o *stringArray) Set(value string) error {
	*o = append(*o, value)
	return nil
}

type cmdRender struct {
	Cmd            *flag.FlagSet
	Overrides      stringArray
	ValuesJsonPath *string
	OutPath        *string
	BlueprintPath  *string
	NoInput        *bool
	ForceInput     *bool
	Help           *bool
}

func CmdRender() *cmdRender {
	command := flag.NewFlagSet("render", flag.ExitOnError)

	cmd := cmdRender{Cmd: command, Overrides: stringArray{}}

	command.Var(&cmd.Overrides, "set", `Set arguments overrides in format "arg=value". Repeat for setting multiple arguments values.`)
	cmd.ValuesJsonPath = command.String("values", "", `Path to arguments values JSON file.`)
	cmd.OutPath = command.String("out", ".", `Path to output rendered template.`)
	cmd.BlueprintPath = command.String("blueprint", "rendr.yaml", `Path to blueprint file inside of template.`)
	cmd.NoInput = command.Bool("noinput", false, `Do not request user input for missing arguments values.`)
	cmd.ForceInput = command.Bool("forceinput", false, `Force user input requests even for noinput arguments.`)
	cmd.Help = command.Bool("help", false, `Prints command help.`)

	command.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: rendr [options] <template-url> [<extra-root> ...]\n")
		fmt.Println()
		fmt.Println("Parameters:")
		fmt.Println(`  <template-url>`)
		fmt.Println(`        Location of the template to be rendered.`)
		fmt.Println(`          For git repositories use the full url, for example: "https://github.com/thecompany/therepo.git".`)
		fmt.Println(`          Optionally if a subfloder of the repository is needed it could be specified after the repo url: "https://github.com/thecompany/therepo.git/subfolder".`)
		fmt.Println(`          Local filesystem template could be used via file URI scheme: "file:///./some/path/template".`)
		fmt.Println(`          Github repositories could be used by their slug which will be examded automatically: "github.com/thecompany/therepo".`)
		fmt.Println(`  <extra-root>`)
		fmt.Println(`        Extra root to add to template files.`)
		fmt.Println(`          Similar to <template-url> might be pointing to the file system or git repo.`)
		fmt.Println()
		fmt.Println("Options:")
		command.PrintDefaults()
		fmt.Println()
		fmt.Println(`To print usage run: rendr help`)
	}
	return &cmd
}

func (command *cmdRender) Execute(arguments []string) {
	command.Cmd.Parse(arguments)

	if *command.Help {
		command.Cmd.Usage()
		os.Exit(0)
	}

	if command.Cmd.NArg() < 1 {
		stderr.Println(`Parameter <template-url> is not provided.`)
		fmt.Println()
		command.Cmd.Usage()
		os.Exit(1)
	}
	templateUrl := command.Cmd.Arg(0)

	extraRoots := []string{}
	for iarg := 1; iarg < command.Cmd.NArg(); iarg++ {
		extraRoot := command.Cmd.Arg(iarg)
		extraRoots = append(extraRoots, normalizeTemplateUrl(extraRoot))
	}

	templateUrl = normalizeTemplateUrl(templateUrl)

	var valuesJsonData []byte = nil
	if *command.ValuesJsonPath != "" {
		data, err := ioutil.ReadFile(*command.ValuesJsonPath)
		failIfError(err, `Failed to read arguments JSON file "%s"`, *command.ValuesJsonPath)
		valuesJsonData = data
	}

	inputMode := render.RegularInputMode
	if *command.ForceInput {
		inputMode = render.ForceInputMode
	}
	if *command.NoInput {
		inputMode = render.NoInputMode
	}

	template := render.Template{templateUrl, *command.BlueprintPath, extraRoots}
	renderedFiles, err := template.Render(inputMode, valuesJsonData, command.Overrides)
	failIfError(err, `Failed to render`)

	err = renderedFiles.WriteAll(*command.OutPath, true)
	failIfError(err, `Failed to write rendered files`)
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
