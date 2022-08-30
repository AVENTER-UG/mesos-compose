# Host Network

The networkmode "host" will start the container with the exposed port inside of the container. 

Example: 

```yaml
version: '3.9'
services:
  redis-host:
    image: redis:latest
    network_mode: "host"
```

`docker ps` will show as the container:

```bash
docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED          STATUS          PORTS                     NAMES
52b853e8272b   redis:latest   "docker-entrypoint.s…"   16 minutes ago   Up 16 minutes                             mesos-16aa16b9-903f-45cc-bd77-5aacfce5d87d
```
