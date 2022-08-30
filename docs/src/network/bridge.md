# Bridged Network

The networkmode "bridge" will start the container with a dynamic host port. 

Example: 

```yaml
version: '3.9'
services:
  redis-bridge:
    image: redis:latest
    ports:
      - "9401:6379"
    network_mode: "bridge"
```

`docker ps` will show as the container:

```bash
docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED          STATUS          PORTS                     NAMES
b15c1ba19cc7   redis:latest   "docker-entrypoint.s…"   11 minutes ago   Up 11 minutes   0.0.0.0:31916->6379/tcp   mesos-42da8556-5bec-4bba-b8e0-98d2e74c6d20]
```
