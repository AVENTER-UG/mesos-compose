# Mesos-Compose CLI Usage

The Mesos-Compose framework does support the new version of mesos-cli.


The following parameters are currently supported:

```bash

mesos-cli compose help
Interacts with the Mesos-Compose Framework

Usage:
  mesos compose (-h | --help)
  mesos compose --version
  mesos compose <command> (-h | --help)
  mesos compose [options] <command> [<args>...]

Options:
  -h --help  Show this screen.
  --version  Show version info.

Commands:
  info     Get information about the running Mesos compose framework.
  kill     Kill Mesos compose workload
  launch   Launch Mesos workload from compose file
  list     Show all running tasks.
  restart  Restart service
  update   Update service from compose file
  version  Get the version number of Mesos compose

```

## Launch Workload

To launch workload, you need a compose file.


```bash

mesos-cli compose launch
Launch Mesos workload from compose file

Usage:
  mesos compose launch (-h | --help)
  mesos compose launch --version
  mesos compose launch [options] <framework-name> <project> <compose-file>

Options:
  -h --help  Show this screen.

Description:
  Launch Mesos workload from compose file
```

Example:

```bash

mesos compose launch mc allwebserver docs/example/test-http.yaml 

```

- `mc` is the Mesos registration name of the framework.
- `allwebserver` is the project name. We can also see it as subcategory. 

## Update Workload

To launch workload, you need a compose file.


```bash

mesos-cli compose update
Update service from compose file

Usage:
  mesos compose update (-h | --help)
  mesos compose update --version
  mesos compose update [options] <framework-name> <project> <compose-file>

Options:
  -h --help  Show this screen.

Description:
  Update service from compose file
```

Example:

```bash

mesos compose update mc allwebserver test1 docs/example/test-http.yaml 

```

- `mc` is the Mesos registration name of the framework.
- `allwebserver` is the project name. We can also see it as subcategory. 
- `test1` is the service name of the container we defined in our compose file.

## Restart Workload

To launch workload, you need a compose file.


```bash

mesos-cli compose restart
Restart service

Usage:
  mesos compose restart (-h | --help)
  mesos compose restart --version
  mesos compose restart [options] <framework-name> <project> <service>

Options:
  -h --help  Show this screen.

Description:
  Restart service
```

Example:

```bash

mesos compose restart mc allwebserver test1 

```

- `mc` is the Mesos registration name of the framework.
- `allwebserver` is the project name. We can also see it as subcategory. 
- `test1` is the service name of the container we defined in our compose file.

## Kill Workload

Kill running or staled workload managed by Mesos-Compose.


```bash

mesos compose kill
Kill Mesos compose workload

Usage:
  mesos compose kill (-h | --help)
  mesos compose kill --version
  mesos compose kill [options] <framework-name> <project> <service> 

Options:
  -h --help  Show this screen.

Description:
  Kill Mesos compose workload

```

Example:

```bash

mesos compose kill mc allwebserver test1

```

- `mc` is the Mesos registration name of the framework.
- `allwebserver` is the project name. We can also see it as subcategory. 
- `test1` is the service name of the container we defined in our compose file.
