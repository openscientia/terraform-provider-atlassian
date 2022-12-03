# tfwaff

`tfwaff` is a Terraform ATLASSIAN Provider scaffolding CLI tool that generates the required files (including acceptance test files) to implement new resources or data sources following best practices and conventions.

## Install

After cloning the repository, navigate to the `tfwaff` directory.

Then, execute the following command to install the application:

```console
go install .
```

`tfwaff` will be installed in your machine's `$GOPATH`.

## Run

Go to the [`internal/provider`](../internal/provider/) directory where all resources and data sources files are located.

To generate resource files:

```console
tfwaff resource --name JiraIssueField

OR 

tfwaff resource --n JiraIssueField
```

To generate data source files:

```console
tfwaff datasource --name JiraIssueField

OR

tfwaff datasource -n JiraIssueField
```

## Commands

### Help

```console
tfwaff is a CLI application that generates the required files
to implement new resources and data sources in the Terraform ATLASSIAN Provider.

Usage:
  tfwaff [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  datasource  Generate all necessary files for a data source
  help        Help about any command
  resource    Generate all necessary files for a resource

Flags:
      --dry-run   do not create or overwrite files
  -f, --force     force creation, overwrite existing files
  -h, --help      help for tfwaff

Use "tfwaff [command] --help" for more information about a command.
```

### Autocompletion

```console
tfwaff completion -h
Generate the autocompletion script for tfwaff for the specified shell.
See each sub-command's help for details on how to use the generated script.

Usage:
  tfwaff completion [command]

Available Commands:
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh

Flags:
  -h, --help   help for completion

Global Flags:
      --dry-run   do not create or overwrite files
  -f, --force     force creation, overwrite existing files

Use "tfwaff completion [command] --help" for more information about a command.
```

### Resource

```console
tfwaff resource -h
Generate all necessary files for a resource

Usage:
  tfwaff resource [flags]

Flags:
  -h, --help          help for resource
  -n, --name string   Name of the new resource in pascal case (i.e. MixedMaps) as: <Service><Name>

Global Flags:
      --dry-run   do not create or overwrite files
  -f, --force     force creation, overwrite existing files
```

### Data Source

```console
tfwaff datasource -h
Generate all necessary files for a data source

Usage:
  tfwaff datasource [flags]

Flags:
  -h, --help          help for datasource
  -n, --name string   Full name of the new data-source in snake case, e.g. <provider>_<service>_<name>

Global Flags:
      --dry-run   do not create or overwrite files
  -f, --force     force creation, overwrite existing files
```
