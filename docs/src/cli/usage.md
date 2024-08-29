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
    framework  Framework Commands.
    info       Get information about the running Mesos compose framework.
    kill       Kill a single task (ID) or a whole service (Task Name)
    launch     Launch Mesos workload from compose file
    list       Show all running tasks.
    restart    Restart a single task (ID) or a whole service (Task Name)
    update     Update service from compose file
    version    Get the version number of Mesos compose

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

## List all Tasks managed my the framework

```bash

mesos-cli compose list mc
ID                                                 Task Name             State         Mesos Agent
test_test1.55662bcc-7268-905e-333a-47a03314d7d5.0  mc:test:allwebserver  TASK_RUNNING  testagent.test.internal

```

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

## Restart all tasks of a service

A Service can run multiple instances of a task. The following example will show,
how to restart the entire service.


```bash

mesos-cli compose restart
Restart a single task (ID) or a whole service (Task Name)

Usage:
  mesos compose restart (-h | --help)
  mesos compose restart --version
  mesos compose restart [options] <framework-name> <task>

Options:
  -h --help  Show this screen.

Description:
  Use the "ID" to restart a single task or the "Task Name" to restart the entire service.

```

Example:

```bash

mesos compose restart mc mc:test:allwebserver

```

- `mc` is the Mesos registration name of the framework.
- `mc:test:allwebserver` is the Task-Name.

## Restart a single Task

Sometimes it's enough to restart a single task and not the entire service.

```bash

mesos-cli compose restart
Restart a single task (ID) or a whole service (Task Name)

Usage:
  mesos compose restart (-h | --help)
  mesos compose restart --version
  mesos compose restart [options] <framework-name> <task>

Options:
  -h --help  Show this screen.

Description:
  Use the "ID" to restart a single task or the "Task Name" to restart the entire service.

```

Example:

```bash

mesos compose restart mc test_test1.55662bcc-7268-905e-333a-47a03314d7d5.0

```

- `mc` is the Mesos registration name of the framework.
- `test_test1.55662bcc-7268-905e-333a-47a03314d7d5.0` is the ID of the Task we want restart.

## Kill a Service or a single Task

To kill a service or a single task is equvalent to restart.

```bash

mesos-cli compose kill
Kill a single task (ID) or a whole service (Task Name)

Usage:
  mesos compose kill (-h | --help)
  mesos compose kill --version
  mesos compose kill [options] <framework-name> <task>

Options:
  -h --help  Show this screen.

Description:
  Use the "ID" to Kill a single task or the "Task Name" to kill the entire service.

```
