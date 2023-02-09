# Contributing

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

- install git pre-commit tool global(macOS)

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

```sh
git-chglog -o CHANGELOG.md
```

### Find ignored files

```sh
find . -type f  | git check-ignore --stdin
```
