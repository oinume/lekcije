[![CircleCI](https://circleci.com/gh/oinume/lekcije/tree/master.svg?style=svg)](https://circleci.com/gh/oinume/lekcije/tree/master)
[![codecov](https://codecov.io/gh/oinume/lekcije/branch/master/graph/badge.svg)](https://codecov.io/gh/oinume/lekcije)

# lekcije
Follow your favorite teachers in DMM Eikaiwa and receive notification when favorite teachers open lessons.

# Install dependencies

### Install docker

```bash
brew cask install dockertoolbox
```

OR 

```bash
brew cask install docker
```

### Install other tools

```bash
brew install dep fswatch node
make setup
npm install
```

### For developers
```bash
brew install chromedriver
```

# Develop on your local machine

### Run MySQL server

with docker-machine
```
docker-machine start default
eval "$(docker-machine env default)"
docker-compose up
```

OR with `Docker for Mac`.
```
docker-compose up
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

OR just use `source` command.
```
source .envrc
```

### Migrate DB
```
goose -env=local up
```

### Run server
```
make watch
```

### Run frontend dev server
```
npm start
```

### Access to the web

http://localhost:4000/

### Connect to MySQL on Docker

Use `docker-machine ip default` on docker-machine
```
mysql -uroot -proot -h $(docker-machine ip default) -P 13306 lekcije
```

OR `127.0.0.1` on Docker for Mac.

```
mysql -uroot -proot -h 127.0.0.1 -P 13306 lekcije
```
