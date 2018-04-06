[![CircleCI](https://circleci.com/gh/oinume/lekcije/tree/master.svg?style=svg)](https://circleci.com/gh/oinume/lekcije/tree/master)
[![codecov](https://codecov.io/gh/oinume/lekcije/branch/master/graph/badge.svg)](https://codecov.io/gh/oinume/lekcije)

# lekcije
Follow your favorite teachers in DMM Eikaiwa and receive notification when favorite teachers open lessons.

# Install dependencies

```bash
brew cask install dockertoolbox
brew install dep fswatch
make setup
npm install
```

## For developers
```bash
brew install chromedriver
```

# Develop on your local machine

### Run MySQL server

```
docker-machine start default
eval "$(docker-machine env default)"
docker-compose up
```

### Migrate DB
```
goose -env=local up
```

### Define environmental variables
```
cp .envrc.local .envrc
```

Use direnv
```
direnv allow
```

or

Just source
```
source .envrc
```

### Run server
```
make watch
```

### Run frontend
```
npm start
```

### Connect to mysql on Docker

```
mysql -uroot -proot -h $(docker-machine ip default) -P 13306 lekcije
```
