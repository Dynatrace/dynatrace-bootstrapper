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
      - Used to create the `<config-directory>/<container-name>/oneagent/agent/config/ruxitagentproc.conf` file
    - `initial-connect-retry`: A file containing a single number value. Defines the delay before the initial connection attempt. (Useful in case of `istio-proxy` is used.)
      - Used to create/update the `<config-directory>/<container-name>/oneagent/agent/customkeys/curl_options.conf` file.
    - `trusted.pem`: A file containing the **certificates** used by the CodeModule for all its communication (proxy communication's not included).
      - Used to create the `<config-directory>/<container-name>/oneagent/agent/customkeys/custom.pem` and `<config-directory>/<container-name>/oneagent/agent/customkeys/custom_proxy.pem`file.
    - `activegate.pem`: A file containing the **certificates** used by the CodeModule for all its communication with the ActiveGate (proxy communication's **NOT** included).
      - Used to create the `<config-directory>/<container-name>/oneagent/agent/customkeys/custom.pem`.
      - Is concatenated with the `trusted.pem` if both is present.
    - `endpoint.properties`: A file containing the necessary info so the metadata-enrichment metrics can be ingested properly
      - Used to create the `<config-directory>/<container-name>/enrichment/endpoint/endpoint.properties`.
      - Example:

      ```conf
      DT_METRICS_INGEST_URL=http://test-activegate.dynatrace/e/tenant/api/v2/metrics/ingest
      DT_METRICS_INGEST_API_TOKEN=<data-ingest-token>
      ```

#### `--install-path`

*Example*: `--input-directory="/opt/dynatrace/oneagent"`

- This is an **optional** arg
  - Defaults to `/opt/dynatrace/oneagent`
- The `--install-path` arg defines the base path where the agent binary will be put. This is only necessary to properly configure the `ld.so.preload` file.
  - The `ld.so.preload` is put under `<config-directory>/oneagent/ld.so.preload`

#### `--fullstack`

*Example*: `--fullstack`

- This is an **optional** arg
  - Defaults to `false`
- The `--fullstack` arg will make sure that the CodeModule is configured to be in fullstack mode.
  - Adds additional values to the `<config-directory>/<container-name>/oneagent/agent/config/container.conf`.

#### `--tenant`

*Example*: `--tenant="my-tenant"`

- This is an **optional** arg, but mandatory incase of `--fullstack`.
- Only used incase of `--fullstack`, provides additional info needed to properly configure `<config-directory>/<container-name>/oneagent/agent/config/container.conf`.

#### `--attribute`

*Example*: `--attribute="k8s.pod.name=test"`

- This is an **optional** arg
- The `--attribute` arg defines the passed in Pod attributes that will be used to configure the metadata-enrichment and injected CodeModule. It is a key value pair.

#### `--attribute-container`

*Example*: `--attribute-container="{k8s.container.name=app, container_image.registry=gcr.io, container_image.repository=test}"`

- This is an **optional** arg
- The `--attribute-container` arg defines the passed in Container attributes that will be used to configure the metadata-enrichment and injected CodeModule. It is a JSON formatted string.

#### `--suppress-error`

*Example*: `--suppress-error`

- This is an **optional** arg
  - Defaults to `false`
- The `--suppress-error` arg will silence any errors, causing the executable to return with an exit code 0 even if an error occurred. Intended purpose is to not block the application container from starting, if used as init-container.

#### `--debug`

*Example*: `--debug`

- This is an **optional** arg
  - Defaults to `false`
- The `--debug` arg will enabled the debug logs.

## Development

- To run tests: `make test`
- To run linting: `make lint`
- To build a testing image: `make build`

### Debug logs

- The flag `--debug` enables debug logs.
- Create new debug logs via `log.V(1).Info(...)`
  - The logging levels between `zap` and `logr` is different, so we had to get creative.

### How to test (in K8s)

#### helm-sample

A a Helm chart for a small PHP sample app, with the bootsrapper as it's `initContainer`, as it is expected to be used.

To deploy it in your cluster (where your `KUBECONFIG` is pointing) with the `image` of the current branch (considering that the `image` was built via `make build`):

- `make deploy`

To remove it from your cluster:

- `make undeploy`

##### Custom values

You can use the `make deploy/custom`

- This will use the `hack/testing/helm-sample/_values.yaml` as the values.
  - This file is ignored by git, so you can safely put whatever you want into it.
- The `image` will still be set according to what `make build` would create.

## How to release

The bootstrapper is mainly used by installing it via `go install`, so the release is very simple:

1. Go to the [release page](https://github.com/Dynatrace/dynatrace-bootstrapper/releases).
2. Press [Draft new release](https://github.com/Dynatrace/dynatrace-bootstrapper/releases/new)
3. Set the tag you want to release, you don't need to push the tag beforehand. Example format: v1.2.3
4. Press `Generate release notes`
5. Cleanup generated output to only have the interesting commits. (business logic relevant or bugfixes)
6. Publish release
