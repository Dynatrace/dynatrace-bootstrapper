# dynatrace-bootstrapper

The `dynatrace-bootstrapper` is a small CLI binary built into a [Dynatrace CodeModule](https://gallery.ecr.aws/dynatrace/dynatrace-codemodules) (after version `<to-be-determined-codemodule-version>`) so that it can be used in conjunction with the [dynatrace-operator](https://github.com/Dynatrace/dynatrace-operator) to optimize the download and configuration of a CodeModule in Kubernetes.

## Responsibilities

- Copy a Dynatrace CodeModule to a target directory
- Configure the Dynatrace CodeModule according to the configuration options provided

## How to use

### Args

#### `--source`

*Example*: `--source="/opt/dynatrace/oneagent"`

- ⚠️This is a **required** arg⚠️
- The `--source` arg defines the base path where to copy the CodeModule FROM.

#### `--target`

*Example*: `--target="example/bins/1.2.3"`

- ⚠️This is a **required** arg⚠️
- The `--target` arg defines the base path where to copy the CodeModule TO.

#### `--work`

*Example*: `--work="/example/work"`

- This is an **optional** arg
- The `--work` arg defines the base path for a tmp folder, this is where the command will do its work, to make sure the operations are atomic. It must be on the same disk as the target folder.

#### `--technology`

*Example*: `--technology="python,java"`

- This is an **optional** arg
- The `--technology` arg defines the paths associated to the given technology in the `manifest.json` file. Only those files will be copied that match the technology. It is a comma-separated list.

## Development

- To run tests: `make test`
- To run linting: `make lint`
- To build a testing image: `make build`
