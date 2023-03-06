# TOC MACHINE TRADING

[![BUILD](https://img.shields.io/github/actions/workflow/status/ToC-Taiwan/toc-machine-trading/main.yml?style=for-the-badge&logo=github)](https://github.com/ToC-Taiwan/toc-machine-trading/actions/workflows/main.yml)
[![Go](https://img.shields.io/github/go-mod/go-version/ToC-Taiwan/toc-machine-trading?style=for-the-badge&logo=go)](https://golang.org)

[![CONTAINER](https://img.shields.io/badge/Container-Docker-blue?style=for-the-badge&logo=docker&logoColor=blue)](https://www.docker.com)
[![IMAGE](https://img.shields.io/docker/pulls/maochindada/toc-machine-trading?style=for-the-badge)](https://hub.docker.com/repository/docker/maochindada/toc-machine-trading/general)
[![IMAGE_SIZE](https://img.shields.io/docker/image-size/maochindada/toc-machine-trading/latest?style=for-the-badge)](https://hub.docker.com/repository/docker/maochindada/toc-machine-trading/general)

[![RELEASE](https://img.shields.io/github/release/ToC-Taiwan/toc-machine-trading?style=for-the-badge)](https://github.com/ToC-Taiwan/toc-machine-trading/releases/latest)
[![LICENSE](https://img.shields.io/github/license/ToC-Taiwan/toc-machine-trading?style=for-the-badge)](COPYING)

## Structure

![Example](docs/img/layers.png)

### Config

```sh
cp ./configs/default.config.yml ./configs/config.yml
```

### Env

```sh
cp .env.template .env
```

### Make

- show help

```sh
make help
```

- run

```sh
make
```

### golangci-lint config from gitlab

```sh
docker run -it registry.gitlab.com/gitlab-org/gitlab-build-images:golangci-lint-alpine cat /golangci/.golangci.yml
```

## Authors

- [**Tim Hsu**](https://github.com/Chindada)
