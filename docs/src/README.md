# Introduction - mesos-compose, the docker-compose framework for Apache Mesos

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
    command: ["sleep","10"]
    restart: always
    volumes:
      - /tmp:/tmp
    environment:
      - MYSQL_HOST=test
    labels:
      biz.aventer.mesos_compose.container_type: "DOCKER"
    network_mode: "BRIDGE"
    network:
      - default
    deploy:
      resources:
        limits:
          cpus: "0.001"
          memory: "50"

networks:
  default:
    external: true
    name: weave

```

Push these compose file to the framework. Every compose file needs to have an
own project name.

```bash
curl -X PUT http://localhost:10000/v0/compose/<PROJECTNAME> --data-binary @docs/example/docker-compose.yml
```

![image_2021-11-08-11-33-09](vx_images/image_2021-11-08-11-33-09.png)

![image_2021-11-08-11-33-47](vx_images/image_2021-11-08-11-33-47.png)

To scale the service, just execute the same call again. To update a already existing docker-compose project, call:

```bash
curl -X PUT http://localhost:10000/v0/compose/<PROJECTNAME>/update --data-binary @docs/example/docker-compose.yml
```
