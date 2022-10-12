# TOC MACHINE TRADING

[![Workflow](https://github.com/ToC-Taiwan/toc-machine-trading/actions/workflows/main.yml/badge.svg)](https://github.com/ToC-Taiwan/toc-machine-trading/actions/workflows/main.yml)
[![Maintained](https://img.shields.io/badge/Maintained-yes-green)](https://github.com/ToC-Taiwan/toc-sinopac-python)
[![Go](https://img.shields.io/badge/Go-1.19.2-blue?logo=go&logoColor=blue)](https://golang.org)
[![OS](https://img.shields.io/badge/OS-Linux-orange?logo=linux&logoColor=orange)](https://www.linux.org/)
[![Container](https://img.shields.io/badge/Container-Docker-blue?logo=docker&logoColor=blue)](https://www.docker.com/)

## Layers

![Example](docs/img/layers.png)

## Tools

### Conventional Commit

- install git cz tool global

```sh
sudo npm install -g commitizen
sudo npm install -g cz-conventional-changelog
sudo npm install -g conventional-changelog-cli
echo '{ "path": "cz-conventional-changelog" }' > ~/.czrc
```

### Pre-commit

- install git pre-commit tool global

```sh
brew install pre-commit
```

- install/modify from config

```sh
pre-commit autoupdate
pre-commit install
pre-commit run --all-files
```

### Modify CHANGELOG

- First Time

```sh
conventional-changelog -p angular -i CHANGELOG.md -s -r 0
```

- From Last semver tag

```sh
conventional-changelog -p angular -i CHANGELOG.md -s
```

### Find ignored files

```sh
find . -type f  | git check-ignore --stdin
```

### Migrate Tool

- install

```sh
version=v4.15.2
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$version
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

### env file example

```sh
echo 'DEPLOYMENT=dev
LOG_FORMAT=console
LOG_LEVEL=trace
DISABLE_SWAGGER_HTTP_HANDLER=
GIN_MODE=release
SINOPAC_URL=127.0.0.1:56666
PG_URL=postgres://postgres:asdf0000@127.0.0.1:5432/
RABBITMQ_URL=amqp://admin:password@127.0.0.1:5672/%2f
RABBITMQ_EXCHANGE=toc
DB_NAME=machine_trade
' > .env
```

### VSCode Debug Setting

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug",
            "type": "go",
            "request": "attach",
            "debugAdapter": "dlv-dap",
            "processId": "toc-machine-trading",
        }
    ]
}
```

### Config

```sh
cp ./configs/default.config.yml ./configs/config.yml
```

## Authors

- [**Tim Hsu**](https://github.com/Chindada)
