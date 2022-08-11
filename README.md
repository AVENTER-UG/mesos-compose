# mesos-compose

[![Chat](https://img.shields.io/static/v1?label=Chat&message=Support&color=brightgreen)](https://matrix.to/#/#mesosk3s:matrix.aventer.biz?via=matrix.aventer.biz)
[![Docs](https://img.shields.io/static/v1?label=Docs&message=Support&color=brightgreen)](https://aventer-ug.github.io/mesos-m3s/index.html)

Mesos Framework to use docker-compose files.

## Requirements

- Apache Mesos min 1.6.0
- Mesos with SSL and Authentication is optional
- Redis Database
- Docker Compose Spec 3.9

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
