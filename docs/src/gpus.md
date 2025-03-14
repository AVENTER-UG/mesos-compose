# GPUs

mesos-compose does support the usage of GPU's with the docker containerizer.
To do so there, take a look into the following examples:

## AMD

If you want to use your AMD GPU's inside of your docker container, you have
to add the following yaml into your mesos-compose.yaml

```yaml
services:
  app:
    ...
    gpus:
      driver: "amd"
```

That will add the following parameters to the docker executor:

```bash
--device=/dev/kfd
--device=/dev/dri
--security-opt seccomp=unconfined
```

## NVIDIA

If you want to use your NVIDIA GPU's inside of your docker container, you have
to add the following yaml into your mesos-compose.yaml

```yaml
services:
  app:
    ...
    gpus:
      driver: "nvidia"
      device: 1
```

That will add the following parameters to the docker executor:

```bash
--gpus device=1
```

The device is the ID or number of your GPU.

