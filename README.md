[![Actions/backend](https://github.com/oinume/lekcije/workflows/backend/badge.svg?branch=master)](https://github.com/oinume/lekcije/actions?query=workflow%3Abackend+branch%3Amaster)
[![Actions/frontend](https://github.com/oinume/lekcije/workflows/frontend/badge.svg?branch=master)](https://github.com/oinume/lekcije/actions?query=workflow%3Abackend+branch%3Amaster)
[![codecov/backend](https://codecov.io/gh/oinume/lekcije/branch/main/graph/badge.svg?flag=backend)](https://codecov.io/gh/oinume/lekcije?flag=backend)
[![codecov/frontend](https://codecov.io/gh/oinume/lekcije/branch/main/graph/badge.svg?flag=frontend)](https://codecov.io/gh/oinume/lekcije?flag=frontend)

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
brew install go fswatch node
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
cp .env.sample .env
```

Replace `192.168.99.100` to `127.0.0.1` on your .env if you use `Docker for Mac`.

And then, load environmental variables with `direnv`.

```
direnv allow
```

OR just use `source` command.
```
source .env
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


