---
title: wp-gitea-release
---

[![Build Status](https://ci.thegeeklab.de/api/badges/thegeeklab/wp-gitea-release/status.svg)](https://ci.thegeeklab.de/repos/thegeeklab/wp-gitea-release)
[![Docker Hub](https://img.shields.io/badge/dockerhub-latest-blue.svg?logo=docker&logoColor=white)](https://hub.docker.com/r/thegeeklab/wp-gitea-release)
[![Quay.io](https://img.shields.io/badge/quay-latest-blue.svg?logo=docker&logoColor=white)](https://quay.io/repository/thegeeklab/wp-gitea-release)
[![Go Report Card](https://goreportcard.com/badge/github.com/thegeeklab/wp-gitea-release)](https://goreportcard.com/report/github.com/thegeeklab/wp-gitea-release)
[![GitHub contributors](https://img.shields.io/github/contributors/thegeeklab/wp-gitea-release)](https://github.com/thegeeklab/wp-gitea-release/graphs/contributors)
[![Source: GitHub](https://img.shields.io/badge/source-github-blue.svg?logo=github&logoColor=white)](https://github.com/thegeeklab/wp-gitea-release)
[![License: Apache-2.0](https://img.shields.io/github/license/thegeeklab/wp-gitea-release)](https://github.com/thegeeklab/wp-gitea-release/blob/main/LICENSE)

Woodpecker CI plugin to add comments to GitHub Issues and Pull Requests.

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< toc >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Usage

{{< hint type=important >}}
Due to the nature of this plugin, a secret for the GitHub token may need to be exposed for pull request events in Woodpecker. Please be careful with this option, as a malicious actor may submit a pull request that exposes your secrets. Do not disclose secrets to pull requests in public environments without further protection.
{{< /hint >}}

{{< hint type=note >}}
Only pull request events are supported by this plugin. Running the plugin on other events will result in an error.
{{< /hint >}}

```YAML
kind: pipeline
name: default

steps:
  - name: pr-comment
    image: quay.io/thegeeklab/wp-gitea-release
    settings:
      api_key: ghp_3LbMg9Kncpdkhjp3bh3dMnKNXLjVMTsXk4sM
      message: "CI run completed successfully"
      update: true
```

### Parameters

<!-- prettier-ignore-start -->
<!-- spellchecker-disable -->
{{< propertylist name=wp-gitea-release.data sort=name >}}
<!-- spellchecker-enable -->
<!-- prettier-ignore-end -->

## Build

Build the binary with the following command:

```shell
make build
```

Build the container image with the following command:

```shell
docker build --file Containerfile.multiarch --tag thegeeklab/wp-gitea-release .
```

## Test

```Shell
docker run --rm \
  -e CI_PIPELINE_EVENT=pull_request \
  -e CI_REPO_OWNER=octocat \
  -e CI_REPO_NAME=foo \
  -e CI_COMMIT_PULL_REQUEST=1
  -e PLUGIN_API_KEY=abc123 \
  -e PLUGIN_MESSAGE="Demo comment" \
  -v $(pwd):/build:z \
  -w /build \
  thegeeklab/wp-gitea-release
```
