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

OR Use `Docker for Mac`.
```
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

Replace `192.168.99.100` to `127.0.0.1` on your .envrc if you use `Docker for Mac`.

And then, load environmental variables with `direnv`.

```
direnv allow
```

Or just use source.
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

### Connect to MySQL on Docker

Use `docker-machine ip default` on docker-machine
```
mysql -uroot -proot -h $(docker-machine ip default) -P 13306 lekcije
```

Or `127.0.0.1` on Docker for Mac.

```
mysql -uroot -proot -h 127.0.0.1 -P 13306 lekcije
```
