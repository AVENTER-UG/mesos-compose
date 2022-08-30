# User Network

The networkmode "user" will start the container without a port. The container will be deployed in the defined CNI network.

Example: 

```yaml
version: '3.9'
services:
  redis-user:
    image: redis:latest
    ports:
      - "9401:6379"
    network_mode: "user"
    network: default

networks:
  default:
    external: true
    name: mini
```

`docker ps` will show as the container:

```bash
docker ps
CONTAINER ID   IMAGE          COMMAND                  CREATED          STATUS          PORTS                     NAMES
ed0350dd8632   redis:latest   "docker-entrypoint.s…"   16 minutes ago   Up 16 minutes   6379/tcp                  mesos-672d1601-42c6-4a6f-be92-4ddc6f271c55
```
