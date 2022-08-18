# rendr

The really usable templating tool:

* Rendering templates stored on the local file system or in the Github repository.
* Templating syntax based on [Mustache](https://mustache.github.io/) with some enhancements.
* Powerful arguments system integrated in Mustache.
* Multiple ways to supply arguments: command line args, JSON file or user console input.
* Could be used as command line tool or as a library built into your Go program.

## Why?

The rendr tool could be used for creating similar projects from parametrized templates.

Often organizations have a need to unify similar projects abd that could be achieved with rendr.

Any framework requires initial project setup, such a setup could be automated via templates.
Many build tools provide ability to create project from templates (like `npm new` or `sbt new`) however they are tech-stack specific and do not exist for all tech ctacks.

The rendr tool provides convenince for both automation and developers via utilizing command line arguments and user console input.

## Table of Contents

<!-- TOC -->
* [Usage](#usage)
* [Template](#template)
  * [Mustache](#mustache)
  * [Blueprint File](#blueprint-file)
  * [Arguments](#arguments)
    * [Arguments Values](#arguments-values)
    * [Array Arguments](#array-arguments)
    * [Arguments Groups](#arguments-groups)
    * [No Input Arguments](#no-input-arguments)
    * [Arguments in Templates](#arguments-in-templates)
    * [Arguments in Paths](#arguments-in-paths)
  * [Additional Blueprint Features](#additional-blueprint-features)
    * [Rename](#rename)
    * [Executables](#executables)
* [Rendr Command Line](#rendr-command-line)
  * [Installation](#installation)
  * [Arguments via Input](#arguments-via-input)
  * [Arguments via JSON](#arguments-via-json)
  * [Arguments via Command Line](#arguments-via-command-line)
  * [Blueprint Location](#blueprint-location)
  * [Output Location](#output-location)
* [Rendr as a Library](#rendr-as-a-library)
<!-- TOC -->

## Usage

This command will render the template located in [examples/simple](https://github.com/specgen-io/rendr/tree/main/examples/simple) folder of the [github.com/specgen-io/rendr](http://github.com/specgen-io/rendr]) repository.
```bash
rendr render github.com/specgen-io/rendr examples/simple
```

Template could be sourced from the local file system:

```bash
rendr render file:///some_path/rendr-example  # note that the path after file:/// should exist 
```

You can find more about command line interface in [Rendr Command Line](#rendr-command-line) section.

## Template

The rendr template consists of files and folders.
It should have a template blueprint file with essential information and arguments definitions.

### Mustache

Rendr is using [Mustache](https://mustache.github.io/) logic-less template syntax with some extensions.
Refer [Mustache documentation](https://mustache.github.io/mustache.5.html) for syntax details.

Here's an example of template file (java code in this case):

```java
package {{package.value}};

public class {{mainclass.value}} {
    public static void main(String[] args) {
    }
}
```

The `{{package.value}}` and `{{mainclass.value}}` are Mustache references to template arguments which have to be defined in the template blueprint file.
Sections below will discuss how to define such arguments.

### Blueprint File

By default, the blueprint file has to be named `rendr.yaml` and located at the root of the template.

Minimal blueprint file:
```yaml
rendr: 0                  # version of the blueprint format
name: example             # technical name of the template
title: Example template   # human-readable title of the template
```

### Arguments

Blueprint file might have some arguments defined in the `args` field.

Blueprint:
```yaml
rendr: 0
name: example
title: Example template

args:                     # templates arguments
  foo:
    type: string
  bar:
    type: boolean
```

Arguments can have `description` which is used for user input whenever argument should be provided by the user.
The `default` value might be set to the argument and this value will be used in case if override value is not provided.

Blueprint:
```yaml
args:
  foo:
    type: string
    description: the foo          # human-readable description of the argument
    default: foo value            # whenever value is not provided this value will be used for the argument
  bar:
    type: boolean
    description: the bar
    default: yes
```

The `string` and `boolean` arguments values could be referenced via `.value` member.
For the example above the following values are available in the templates: `foo.value` and `bar.value`.

#### Arguments Values

String arguments could have set of `values` to limit what values are allowed for the argument.

Blueprint:
```yaml
args:
  baz:
    type: boolean
    description: the baz
    values: [blip, blop, clunk]   # only these values could be set as baz argument value
```

If `values` are set for the argument then additional members are available in the template besides `.value` (which has the raw value for the argument).
In the example above additional boolean tags are: `baz.blip`, `baz.blop`, `baz.clunk`.
They will indicate if the corresponding value is set or not.

Read [Arguments in Templates](#arguments-in-templates) section for information on using arguments in templates.

#### Array Arguments

Arguments might have ann array (of strings) type. In this case value of this arg would be an array of values.

Blueprint:
```yaml
args:
  baz:
    type: array
    description: the baz
    values: [blip, blop, clunk]
    default: [blip, clunk]
```

In the example above `baz.value` is an array of values. Similarly to string arguments following boolean tags are also available: `baz.blip`, `baz.blop`, `baz.clunk`.

#### Arguments Groups

Arguments could be united into groups for convenience. In the example below version is a group of arguments.

Blueprint:
```yaml
args:
  versions:
    type: group
    args:
      foo:
        type: string
        default: 1.0.0
      bar:
        type: string
        default: 2.0.0
      baz:
        type: string
        default: 3.0.0
```

#### No Input Arguments

All arguments including groups might have `noinput` setting. By default it's `false`.
If `noinput` setting is set to `true` that means: do not request user input for the argument.

Blueprint:
```yaml
args:
  versions:
    type: group
    noinput: true   # do not request versions from user
    args:
      foo:
        type: string
        default: 1.0.0
      bar:
        type: string
        default: 2.0.0
      baz:
        type: string
        default: 3.0.0
```

Read [Arguments Input](#arguments-input) section for more information about user input.

#### Arguments in Templates

Normal Mustache tags substitution works in templates.

Each argument when passed to the Mustache template has `value` member with the raw value that is provided for the arguments.
In the example below values for string arguments `package` and `mainclass` are used in the template.

Blueprint:
```yaml
args:
  package:
    type: string
  mainclass:
    type: string
```
Usage:
```java
package {{package.value}};
//        ^ using package argument

public class {{mainclass.value}} {
//             ^ using mainclass argument
    public static void main(String[] args) {
    }
}
```

Similarly `boolean` arguments have `value` with the `boolean` value: 

Blueprint:
```yaml
args:
  helloworld:
    type: boolean
```
Usage:
```java
package com.example;

public class Main {
    public static void main(String[] args) {
        {{#helloworld.value}}
        System.out.println("Hello world!!!");
        // ^ this line will be rendered only if helloworld argument is true
        {{/helloworld.value}}
    }
}
```

Arguments that have possible values set also populated with boolean flags to enable checks for specific values:

Blueprint:
```yaml
args:
  features:
    type: array
    values: [helloworld, exit]
```
Usage:
```java
package com.example;

public class Main {
    public static void main(String[] args) {
        {{#features.helloworld}}
        System.out.println("Hello world!!!");
        // ^ this will be rendered only if feature helloworld is set
        {{/features.helloworld}}
        {{#features.exit}}
        System.exit(0);
        // ^ this will be rendered only if feature exit is set
        {{/features.exit}}
    }
}
```

#### Arguments in Paths

Mustache syntax could be used in file and folder names.

Arguments can be used in the names of files and folders via same Mustache template syntax:

Blueprint:
```yaml
args:
  main:
    type: string
```

Files:
```
/util.java
/{{main.value}}.java
```

Mustache syntax for sections is too "wordy": `{{#condition}}name{{/condition}}`
So rendr offers shorter syntax only for file and folder names.
Closing tag is optional when arguments are used in file and folder names.

Blueprint:

```yaml
args:
  build:
    type: string
    values: [maven, gradle]
```
Canonical Mustache syntax:
```
/{{#build.maven}}pom.xml{{/build.maven}}          # this included only if argument build has maven value
/{{#build.gradle}}build.gradle{{/build.gradle}}   # this included only if argument build has gradle value
```
Short syntax:
```
/{{#build.maven}}pom.xml
/{{#build.gradle}}build.gradle
```

If for whatever reason the name of folder is rendered to empty string then such folder content is just inlined into parent folder.
The example above could be designed using "empty" folder names:
Short syntax:
```
/{{#build.maven}}
    /pom.xml         # this included only if argument build has maven value
/{{#build.gradle}}
    /build.gradle    # this included only if argument build has maven value
```

### Additional Blueprint Features

#### Rename

The blueprint has an option to rename files while rendering the template.
In the example below the `gitignore` (no leading dot `.`) file in the template is renamed into `.gitignore`.
This is useful because otherwise file `.gitignore` whould be treated as git ignoring instruction for the template itself.

Blueprint:
```yaml
rendr: 0
name: example
title: Example template
rename:
  "gitignore": ".gitignore"
  # more renames could be set here
```

#### Executables

Some scripts or files might be marked as executables during the template rendering.
The example below marks maven wrapper script `mvnw` as executable.

Blueprint:
```yaml
rendr: 0
name: example
title: Example template
executables:
  - "mvnw"
  # more executables could be added here
```


## Rendr Command Line

Rendr command line tool renders template from a local file system or Github repository.

### Installation

The rendr could be installed with `go install` command:

```bash
go install github.com/specgen-io/rendr@<version>
```

Alternatively rendr binary could be downloaded from [repository releases](https://github.com/specgen-io/rendr/releases).

Here's a simple usage example:
```bash
rendr render github.com/specgen-io/rendr examples/simple
#            ^ repo with template        ^ optional path to the template inside of the repo
```

### Arguments via Input

Whenever rendr doesn't have a value for the argument it will request for the user input.
This behaviour could be adjusted with flags `-noinput` and `-forceinput`.

The `-noinput` flag disables user input completely. If the flag is set and there's no value for the specific argument the rendering will fail.
This mode is very useful for automation where user input is not possible.

The `-forceinput` flag forces user input even for those arguments that are marked as `noinput` (check [No Input Arguments](#no-input-arguments) section).

### Arguments via JSON

Arguments values might be provided via JSON file.
This might be useful in automation use cases when dealing with many arguments.

Blueprint:
```yaml
args:
  foo:
    type: string
  bar:
    type: boolean
  versions:
    type: group
    args:
      foo:
        type: string
      bar:
        type: string
```

File `values.json`:
```json
{
  "foo": "the foo value",
  "bar": true,
  "versions": {
    "foo": "3.0.0",
    "bar": "4.0.0"
  }
}
```

Command:
```bash
rendr render -values values.json github.com/specgen-io/rendr examples/simple
#            ^ pass JSON file with arguments values
```

### Arguments via Command Line

Arguments values might be provided via command line.

Blueprint:
```yaml
args:
  foo:
    type: string
  bar:
    type: boolean
  versions:
    type: group
    args:
      foo:
        type: string
      bar:
        type: string
```

Command:
```bash
rendr render -set foo="the foo" -set bar="the bar" -set versions.foo="1.0" -set versions.bar="2.0" github.com/specgen-io/rendr examples/simple
#            ^ set foo argument value              ^ set option can be used multiple times
```

Note how grouped arguments are set by their full names: `versions.foo` and `versions.bar`.

### Blueprint Location

The default location of the blueprint file is `./rendr.yaml`.
This could be customized via `-blueprint` option:

```bash
rendr render -blueprint blueprint.yaml github.com/specgen-io/rendr examples/simple
```

### Output Location

The rendr allows to customize output path.
By default rendr tries to write rendered template into the current folder.
This could be customized with `-out` option:

```bash
rendr render -out ./output/path github.com/specgen-io/rendr examples/simple
```

## Rendr as a Library

Rendr could be used as a library.
This is useful if you want to embed templates rendering into your command line tool for better developer experience. 

Add rendr dependency:
```bash
go get github.com/specgen-io/rendr
```

Here's how template could be rendered:
```go
// get the template
template := render.Template{templateUrl, path, blueprintPath}

// render the template
renderedFiles, err := template.Render(inputMode, valuesJsonData, overrides)

// write files
err = renderedFiles.WriteAll(outPath, true)
```

Check [main.go](https://github.com/specgen-io/rendr/blob/main/main.go) of rendr command line tool to explore working sample code rendering templates.