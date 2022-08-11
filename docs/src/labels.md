# Custom labels

## biz.aventer.mesos_compose.container_type

Values: docker, mesos, none

This label will control, which container engine mesos will choose to execute the
container. The value "none" will tell mesos to run the command with the command-executer.

## biz.aventer.mesos_compose.executor

Values: <mesos-executor>

With this label it is possible to use a custom executor.

## biz.aventer.mesos_compose.executor_uri

Values: <mesos-executor-uri>

The URL where to fetch the executor and other files. As example: 
```
'[{"Value":"https://localhost:1234/executor","OutputFile":"mesos-executor"}]'
```
