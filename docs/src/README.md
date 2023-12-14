# Introduction - mesos-compose, the docker-compose framework for Apache Mesos

## Requirements

- Apache Mesos min 1.6.0
- Mesos with SSL and Authentication is optional
- Redis Database
- Docker Compose Spec 3.9

## Example

Compose file with all supported parameters:

```yaml
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
      traefik.http.routers.test.service: "mc_test_app1_80" # if an alternative taskname is set, we have to use it here to
      traefik.http.routers.test.rule: "HostRegexp(`example.com`, `{subdomain:[a-z]+}.example.com`)"
    network_mode: "BRIDGE"
    ports:
      - "8080:80"
      - "9090"
      - "8081:81/tcp"
      - "8082:82/udp"
      - "8082:82/http"
      - "8082:82/https"  
      - "8082:82/h2c"
      - "8082:82/wss"       
    network: default
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
      tcp:
        port:        
    deploy:
      placement:
        constraints:
          - "node.hostname==localhost"
          - "node.platform.os==linux"
          - "node.platform.arch==arm"
      replicas: 1
      resources:
        limits:
          cpus: 0.01
          memory: 50

networks:
  default:
    external: true
    name: weave
    driver: bridge

volumes:
  12345test:
    driver: local
```

We can also use yaml anchors and more then one service in a compose file:

```bash
version: '3.9'

common: &common
  image: alpine:latest
  restart: always
  labels:
    biz.aventer.mesos_compose.container_type: "DOCKER"

services:
  app1:
    <<: *common
    command: ["sleep", "1000"]

  app2:
    <<: *common
    command: ["sleep", "2000"]

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

To scale the service, just execute the same call again. To update a already existing docker-compose project, call:

```bash
curl -X UPDATE http://localhost:10000/api/compose/v0/<PROJECTNAME> --data-binary @docs/example/docker-compose.yml
```
