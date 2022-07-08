# Custom labels

## biz.aventer.mesos_compose.container_type

Values: docker, mesos, none

This label will control, which container engine mesos will choose to execute the
container. The value "none" will tell mesos to run the command with the command-executer.

## biz.aventer.mesos_compose.contraint_hostname

Values: <hostname>

This label will control, on which node the container will be executed.

## biz.aventer.mesos_compose.executor

Values: <mesos-executor>

With this label it is possible to use a custom executor.

## biz.aventer.mesos_compose.executor_uri

Values: <mesos-executor-uri>

The URL where to fetch the executor. As example: `http://localhost:8080/executor`
