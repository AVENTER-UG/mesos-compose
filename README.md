# mesos-compose

[![Chat](https://img.shields.io/static/v1?label=Chat&message=Support&color=brightgreen)](https://matrix.to/#/#mesosk3s:matrix.aventer.biz?via=matrix.aventer.biz)
[![Docs](https://img.shields.io/static/v1?label=Docs&message=Support&color=brightgreen)](https://aventer-ug.github.io/mesos-m3s/index.html)

Mesos Framework to use docker-compose files.

## Requirements

- Apache Mesos min 1.6.0
- Mesos with SSL and Authentication is optional
- Redis Database

## Run Framework

The following environment parameters are only a example. All parameters and the default values are documented in 
the `init.go` file (real documentation will be coming later). These example assume, that we run mesos-mini.

### Step 1

Run a redis server:

```Bash
docker run --rm --name redis -d -p 6379:6379 redis
```

### Step 2

mesos-compose needs some parameters to connect to Mesos. The following serve only as an example.

```Bash
export MESOS_SSL="false"
export LOGLEVEL="DEBUG"
export DOMAIN=".mini"
export AUTH_USERNAME="user"
export AUTH_PASSWORD="password"
export PORTRANGE_FROM=31000
export PORTRANGE_TO=32000
export SKIP_SSL=true
```

### Step 3

Before we launch mesos-compose, we create dedicated network in docker.

```Bash
docker network create --subnet 10.40.0.0/24 mini
```

### Step 4

Now mesos-compose can be started:

```Bash
./mesos-compose
```

### Mesos-Compose in real Apache Mesos environments

In real mesos environments, we have to set at least the following environment variables:

```Bash
export MESOS_MASTER="leader.mesos:5050"
export MESOS_USERNAME=""
export MESOS_PASSWORD=""
```

Also the following could be usefull.

```Bash
export REDIS_SERVER="127.0.0.1:6379"
export REDIS_PASSWORD=""
export REDIS_DB="1"
export MESOS_CNI="weave"
```

## Example

The compose file:

```bash
version: '3.9'

services:
  app:
    image: alpine:latest
    command: ["sleep", "1000"]
    restart: always
    volumes:
      - "12345test:/tmp"
    environment:
      - MYSQL_HOST=test
    hostname: test
    container_name: test
    labels:
      biz.aventer.mesos_compose.container_type: "DOCKER"
      iz.aventer.mesos_compose.executor: "./my-custom-executor"
      biz.aventer.mesos_compose.executor_uri: "http://localhost/my-custom-executor"
      traefik.enable: "true"
      traefik.http.routers.test.entrypoints: "web"
      traefik.http.routers.test.service: "mc_test_app_80"
      traefik.http.routers.test.rule: "HostRegexp(`example.com`, `{subdomain:[a-z]+}.example.com`)"
    network_mode: "BRIDGE"
    ports:
      - "8080:80"
      - "9090"
      - "8081:81/tcp"
      - "8082:82/udp"
    network:
      - default
    deploy:
      placement:
        constraints:
          - "node.hostname==localhost"
      replicas: 1
      resources:
        limits:
          cpus: "0.01"
          memory: "50"

networks:
  default:
    external: true
    name: weave

volumes:
  12345test:
    driver: local

```


Push these compose file to the framework. Every compose file needs to have an
own project name.

```bash
curl -X PUT http://localhost:10000/api/compose/v0/<PROJECTNAME> --data-binary @docs/example/docker-compose.yml
```

![image_2021-11-08-11-33-09](vx_images/image_2021-11-08-11-33-09.png)

![image_2021-11-08-11-33-47](vx_images/image_2021-11-08-11-33-47.png)
