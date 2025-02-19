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
- The `--technology` arg defines the paths associated to the given technology in the `<source>/manifest.json` file. Only those files will be copied that match the technology. It is a comma-separated list.

#### `--config-directory`

*Example*: `--config-directory="/example/config/dir"`

- This is an **optional** arg
- The `--config-directory` arg defines the path where `enrichment` and `config` files will be created in.

#### `--input-directory`

*Example*: `--input-directory="/example/input"`

- This is an **optional** arg
- The `--input-directory` arg defines the base path where configuration files are provided.
  - Intended use is to mount a k8s `Secret` with the correct content in the provided path.
    - So when we are talking about "files", it is just the fields of the `Secret's data`
  - Config files:
    - `ruxitagentproc.json`: A json file containing a response from the `/deployment/installer/agent/processmoduleconfig` endpoint of the Dynatrace Environment(v1) API.
      - This file is **required** if `--input-directory` is defined.
      - Used to update the `<target>/agent/conf/ruxitagentproc.conf` file
    - `initial-connect-retry`: A file containing a single number value. Defines the delay before the initial connection attempt. (Useful in case of `istio-proxy` is used.)
      - Used to create/update the `<config-directory>/<container-name>/oneagent/agent/customkeys/curl_options.conf` file.
    - `trusted.pem`: A file containing the **certificates** used by the CodeModule for all its communication (proxy communication's not included).
      - Used to create the `<config-directory>/<container-name>/oneagent/agent/customkeys/custom.pem` and `<config-directory>/<container-name>/oneagent/agent/customkeys/custom_proxy.pem`file.
    - `activegate.pem`: A file containing the **certificates** used by the CodeModule for all its communication with the ActiveGate (proxy communication's **NOT** included).
      - Used to create the `<config-directory>/<container-name>/oneagent/agent/customkeys/custom.pem`.
      - Is concatenated with the `trusted.pem` if both is present.

#### `--install-path`

*Example*: `--input-directory="/opt/dynatrace/oneagent"`

- This is an **optional** arg
  - Defaults to `/opt/dynatrace/oneagent`
- The `--install-path` arg defines the base path where the agent binary will be put. This is only necessary to properly configure the `ld.so.preload` file.
  - The `ld.so.preload` is put under `<config-directory>/oneagent/ld.so.preload`

#### `--attribute`

*Example*: `--attribute="k8s.pod.name=test"`

- This is an **optional** arg
- The `--attribute` arg defines the passed in Pod attributes that will be used to configure the metadata-enrichment and injected CodeModule. It is a key value pair.

#### `--attribute-container`

*Example*: `--attribute-container="{k8s.container.name=app, container_image.registry=gcr.io, container_image.repository=test}"`

- This is an **optional** arg
- The `--attribute-container` arg defines the passed in Container attributes that will be used to configure the metadata-enrichment and injected CodeModule. It is a JSON formatted string.

## Development

- To run tests: `make test`
- To run linting: `make lint`
- To build a testing image: `make build`

### Debug logs

- The flag `--debug` enables debug logs.
- Create new debug logs via `log.V(1).Info(...)`
  - The logging levels between `zap` and `logr` is different, so we had to get creative.
