package main

import (
	"flag"
	"fmt"
	"github.com/specgen-io/rendr/files"
	"github.com/specgen-io/rendr/render"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var stderr = log.New(os.Stderr, "", 0)

func main() {
	args := parseArguments()

	if strings.HasPrefix(args.TemplateUrl, "github") {
		args.TemplateUrl = fmt.Sprintf(`https://%s.git`, args.TemplateUrl)
	}

	var valuesJsonData []byte = nil
	if args.ValuesJsonPath != "" {
		data, err := ioutil.ReadFile(args.ValuesJsonPath)
		failIfError(err, `Failed to read arguments JSON file "%s"`, args.ValuesJsonPath)
		valuesJsonData = data
	}

	template := render.Template{args.TemplateUrl, args.Path, args.BlueprintPath}
	renderedFiles, err := template.Render(args.OutPath, args.NoInput, args.ForceInput, valuesJsonData, args.Overrides)
	failIfError(err, `Failed to render`)

	err = files.WriteAll(renderedFiles, true)
	failIfError(err, `Failed to write rendered files`)
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

type Arguments struct {
	TemplateUrl    string
	Path           string
	BlueprintPath  string
	ValuesJsonPath string
	Overrides      []string
	OutPath        string
	NoInput        bool
	ForceInput     bool
}

func parseArguments() Arguments {
	var overrides = stringArray{}
	flag.Var(&overrides, "set", `Set arguments overrides in format "arg=value". Repeat for setting multiple arguments values.`)
	valuesJsonPath := flag.String("values", "", `Path to arguments values JSON file.`)
	outPath := flag.String("out", ".", `Path to output rendered template.`)
	blueprintPath := flag.String("blueprint", "blueprint.yaml", `Path to blueprint file inside of template.`)
	noinput := flag.Bool("noinput", false, `Do not request user input for missing arguments values.`)
	forceinput := flag.Bool("forceinput", false, `Force user input requests even for noinput arguments.`)

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: rendr [options] <template-url> [<path>]\n")
		fmt.Println()
		fmt.Println("Parameters:")
		fmt.Println(`  <template-url>`)
		fmt.Println(`        Location of the template to be rendered.`)
		fmt.Println(`          For git repositories use the full url, for example: "https://github.com/thecompany/therepo.git".`)
		fmt.Println(`          Local filesystem template could be used via file URI scheme: "file:///./some/path/template"`)
		fmt.Println(`          Github repositories could be used by their slug: "github.com/thecompany/therepo"`)
		fmt.Println(`  <path>`)
		fmt.Println(`        Path to the root of the template inside of <template-url>. Used only when the repository/folder contains multiple templates. (default "")`)
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println(`To print usage run: rendr help`)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		stderr.Println(`Parameter <template-url> is not provided.`)
		fmt.Println()
		flag.Usage()
		os.Exit(1)
	}

	if flag.Arg(0) == "help" {
		flag.Usage()
		os.Exit(0)
	}

	templateUrl := flag.Arg(0)

	path := ""
	if flag.NArg() > 1 {
		path = flag.Arg(1)
	}

	return Arguments{
		templateUrl,
		path,
		*blueprintPath,
		*valuesJsonPath,
		overrides,
		*outPath,
		*noinput,
		*forceinput,
	}
}
