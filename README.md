# helm-gen-questions

A Helm plugin that generates the [Rancher chart `questions.yml`](https://ranchermanager.docs.rancher.com/how-to-guides/new-user-guides/helm-charts-in-rancher/create-apps#question-variable-reference)
from `x-rancher-question` annotations in a chart's `values.schema.json`, using `values.yaml` for the defaults and types.

The annotations are usually written as `# @schema` comments in `values.yaml` and copied into the schema by the
[helm-schema](https://github.com/dadav/helm-schema) plugin, so run it first.

## Install

Requires Helm 4. The tarball has no GPG `.prov` file, so Helm can't verify it itself and `--verify=false` is needed. Verify the release with the GitHub attestation first (see below):

```sh
helm plugin install --verify=false https://github.com/tashima42/helm-gen-questions/releases/download/vX.Y.Z/gen-questions-X.Y.Z.tgz
```

### Verify the release

Releases are [immutable](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases): once published, the tag and the assets can't be changed or replaced, and GitHub attests them. With the [GitHub CLI](https://cli.github.com) (2.81 or newer):

```sh
# the release and its tag are the ones GitHub attested
gh release verify vX.Y.Z --repo tashima42/helm-gen-questions

# the downloaded tarball matches the attested asset
gh release verify-asset vX.Y.Z gen-questions-X.Y.Z.tgz --repo tashima42/helm-gen-questions
```

Both commands fail if the release was modified or the file doesn't match. The release also has a `checksums.txt` file with the SHA256 of the tarball.

To build from source (Go required):

```sh
git clone https://github.com/tashima42/helm-gen-questions && cd helm-gen-questions
make build
helm plugin install .
```

## Usage

```sh
helm gen-questions [flags]
```

| Flag | Default | Description |
| --- | --- | --- |
| `-v`, `--values-file` | `values.yaml` | Chart values file. `values.yml` is used if the default doesn't exist. |
| `-s`, `--schema-file` | `values.schema.json` | Values JSON schema with the annotations. |
| `-q`, `--questions-file` | `questions.yml` | File to write, `-` for stdout. |
| `--header` | generated file notice | Comment added at the top of the questions file. Each line is prefixed with `# `. An empty value adds none. |
| `--header-file` | | Read the header comment from a file. Can't be combined with `--header`. |

```sh
helm gen-questions -s charts/app/values.schema.json -v charts/app/values.yaml -q charts/app/questions.yml
helm gen-questions --header-file HEADER.txt
helm gen-questions -q - --header ""
```

## Annotations

`x-rancher-question` is an object, or a list of objects when a key needs more than one question.

| Field | |
| --- | --- |
| `label`, `order` | Required. Questions are sorted by `order`, which is not written to the file. |
| `group` | Required, except for subquestions. |
| `variable` | Path relative to the annotated key, needed for keys using `$ref`. |
| `subquestion_of` | Variable of the parent question. Only one level is supported. |
| `type`, `description`, `default` | Override the values inferred from the schema and `values.yaml`. |
| `options`, `empty_option_label` | Enum options, defaulting to the schema `enum`. A `""` option is labeled `None` unless overridden. |

Any other field is copied as is (`required`, `min`, `max`, `show_if`, `show_subquestion_if`...).

## Development

```sh
make test
make build-all   # all the platforms in plugin.yaml
```

Releases are created by pushing a `vX.Y.Z` tag, the plugin version in `plugin.yaml` is set from it.
