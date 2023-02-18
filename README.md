# TOC MACHINE TRADING

[![Workflow](https://github.com/ToC-Taiwan/toc-machine-trading/actions/workflows/main.yml/badge.svg)](https://github.com/ToC-Taiwan/toc-machine-trading/actions/workflows/main.yml)
[![Go](https://img.shields.io/badge/Go-1.20.1-blue?logo=go&logoColor=blue)](https://golang.org)
[![Container](https://img.shields.io/badge/Container-Docker-blue?logo=docker&logoColor=blue)](https://www.docker.com/)
[![Release](https://img.shields.io/github/release/ToC-Taiwan/toc-machine-trading.svg)](https://github.com/golangci/golangci-lint/releases/latest)
[![License](https://img.shields.io/github/license/ToC-Taiwan/toc-machine-trading)](/LICENSE)

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

- build

```sh
make
```

### golangci-lint

```sh
docker run -it registry.gitlab.com/gitlab-org/gitlab-build-images:golangci-lint-alpine bash
find / -name ".golangci.yml"
cat /golangci/.golangci.yml
```

### Dev Note

```log
WARN[2022-09-26T01:23:38+08:00] Fetch Date: 2022-09-18, No Data
WARN[2022-09-26T01:23:41+08:00] Fetch Date: 2022-09-19, FirstTickTime: 2022-09-16 15:00:00, LastTickTime: 2022-09-19 13:44:59, Total: 20037
WARN[2022-09-26T01:23:44+08:00] Fetch Date: 2022-09-20, FirstTickTime: 2022-09-19 15:00:00, LastTickTime: 2022-09-20 13:44:59, Total: 25236
WARN[2022-09-26T01:23:51+08:00] Fetch Date: 2022-09-21, FirstTickTime: 2022-09-20 15:00:00, LastTickTime: 2022-09-21 13:44:59, Total: 43037
WARN[2022-09-26T01:24:19+08:00] Fetch Date: 2022-09-22, FirstTickTime: 2022-09-21 15:00:00, LastTickTime: 2022-09-22 13:44:59, Total: 206679
WARN[2022-09-26T01:24:43+08:00] Fetch Date: 2022-09-23, FirstTickTime: 2022-09-22 15:00:00, LastTickTime: 2022-09-23 13:44:59, Total: 169888
WARN[2022-09-26T01:24:55+08:00] Fetch Date: 2022-09-24, FirstTickTime: 2022-09-23 15:00:00, LastTickTime: 2022-09-24 04:59:59, Total: 87344
WARN[2022-09-26T01:24:55+08:00] Fetch Date: 2022-09-25, No Data
```

## Authors

- [**Tim Hsu**](https://github.com/Chindada)
