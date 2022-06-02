# TOC MACHINE TRADING

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
pre-commit install
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

### Migrate Tool

- install

```sh
version=v4.15.2
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$version
```

### env file example

```sh
echo 'DEPLOYMENT=prod
DISABLE_SWAGGER_HTTP_HANDLER=true
GIN_MODE=release
SINOPAC_URL=127.0.0.1:56666
PG_URL=postgres://user:password@localhost:5432/
DB_NAME=trade' > .env
```

### Config

```sh
cp ./configs/default.config.yml ./configs/config.yml
```

## Authors

- [**Tim Hsu**](https://gitlab.tocraw.com/root)
