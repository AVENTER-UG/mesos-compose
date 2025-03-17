# mesos-compose

[![Docs](https://img.shields.io/static/v1?label=&message=Issues&color=brightgreen)](https://github.com/m3scluster/compose/issues)
[![Chat](https://img.shields.io/static/v1?label=&message=Chat&color=brightgreen)](https://matrix.to/#/#mesosk3s:matrix.aventer.biz?via=matrix.aventer.biz)
[![Docs](https://img.shields.io/static/v1?label=&message=Docs&color=brightgreen)](https://m3scluster.github.io/compose/)
[![GoDoc](https://godoc.org/github.com/AVENTER-UG/mesos-compose?status.svg)](https://godoc.org/github.com/AVENTER-UG/mesos-compose) 
[![Docker Pulls](https://img.shields.io/docker/pulls/avhost/mesos-compose)](https://hub.docker.com/repository/docker/avhost/mesos-compose/)

Mesos Framework to use docker-compose based files.

## Funding

[![](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=H553XE4QJ9GJ8)


## Issues

To open an issue, please use this place: https://github.com/m3scluster/compose/issues

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

```yaml
## This is an example of all supported parameters

version: '3.9'

services:
  app:
    image: alpine:latest
    command: "sleep"
    arguments: ["1000"]
    restart: always
    volumes:
      - "12345test:/tmp"
    environment:
      MYSQL_HOST: test
    hostname: test
    container_name: test
    container_type: docker
    shell: true
    mesos:
      task_name: "mc:test:app1" # an alternative taskname
      executer:
        command: "/mnt/mesos/sandbox/my-custom-executor"
      fetch:
        - value: http://localhost/my-custom-executor
          executable: true
					extract:  false
					cache: false
    labels:
      traefik.enable: "true"
      traefik.http.routers.test.entrypoints: "web"
      traefik.http.routers.test.service: "mc_test_app1_80" # if an alternative taskname is set, we have to use it here too
      traefik.http.routers.test.rule: "HostRegexp(`example.com`, `{subdomain:[a-z]+}.example.com`)"
    network_mode: "BRIDGE" # network mode always overwrite the driver under networks (line 84)
    network: default
    networks:
      default:
        aliases:
          - test # alias hostname
    gpus:
      driver: "nvidia"
      device: 1
    ports:
      - "8080:80"
      - "9090"
      - "8081:81/tcp"
      - "8082:82/udp"
      - "8082:82/http"
      - "8082:82/https"
      - "8082:82/h2c"
      - "8082:82/wss"
    ulimits:
      memlock:
        soft: -1
        hard: -1
      nofile:
        soft: 65536
        hard: 65536
    healthcheck:
      delay_seconds: 15
      interval_seconds: 10
      timeout_seconds: 20
      consecutive_failures: 3
      grace_period_seconds: 10
      command:
        value: "mysqladmin ping -h localhost"
      http:
        scheme:
        port:
        path:
        statuses:
          - 200
      tcp:
        port:
    deploy:
      runtime: "runcvm"
      attributes:
          - "gpus: True"
          - "airflow: True"
      placement:
        constraints:
          - "node.hostname==localhost"
          - "node.platform.os==linux"
          - "node.platform.arch==arm"
          - "unique"
      replicas: 1
      resources:
        limits:
          cpus: 0.01
          memory: 50

networks:
  default:
    external: true
    name: mesos-cni
    driver: user

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
