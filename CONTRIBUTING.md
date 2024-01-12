# Contributing

## Tools

### Conventional Commit

- install git cz tool global

```sh
npm install -g commitizen
npm install -g cz-conventional-changelog
npm install -g conventional-changelog-cli
echo '{ "path": "cz-conventional-changelog" }' > ~/.czrc
```

### Pre-commit

- install pre-commit

```sh
pip3 install pre-commit
```

```sh
pre-commit autoupdate
pre-commit install
```

```sh
pre-commit run --all-files
```

### Modify CHANGELOG

- git-chglog

```sh
brew tap git-chglog/git-chglog
brew install git-chglog
```

- new tag

```sh
COMMIT_HASH=a8d79fe8d00881948a3c42e04fc160a3a90cb4f1
VERSION=2.6.0
git tag -a v$VERSION $COMMIT_HASH -m $VERSION
git-chglog -o CHANGELOG.md
git commit

git push -u origin --tags
git push -u origin --all
```

### Find ignored files

```sh
find . -type f  | git check-ignore --stdin
```
