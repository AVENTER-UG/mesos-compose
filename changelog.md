# Changelog

## v1.2.0

- CHANGE: !!!! If compose.networks is set, the network mode will be automaticly "user".
       It can be overwritten if compose.networks.driver is set. !!!!
- ADD: Enable/Disable event handling as own thread (ENV Variable THREAD_ENABLE)
- ADD: API Endpoint to cleanup the Framework ID. That will force a resubscription under
       a new Framework ID.
- ADD: Support for AMD and NVIDIA GPUs.
- ADD: Mesos Framework GPU Capabilities to get GPU offers from mesos.
- ADD: Runtime support for docker container (https://github.com/newsnowlabs/runcvm)
- ADD: Support for Mesos attributes
- FIX: [API] Restart service
- DEL: [API] Remove useless restart task API
- FIX: Force suppress after successfull framework subscription to prevent unwanted offers
- ADD: GPU Allocation Option in Mesos. Can still use GPU's on the host but will not be allocated in mesos.
- ADD/FIX: TASK_LOST Update causes TASKS to be killed
- FIX: Unwanted Decline Offer causing duplicate declines removed
- ADD: Host constraint: Option to only accept/request offers from particular hosts. 
- FIX: Continously accepting offers when no task running. 

## v1.1.3

- FIX: Do not reconcile not scheduled tasks

## v1.1.2

- FIX: Restart API set the wrong restart flag.
- ADD: Restart API to restart a single Task.
- ADD: Kill API to kill a single Task.
- FIX: Healtcheck missconfiguration
- FIX: Task Lost message handling
- ADD: Wait 60 seconds before try to reconnect to the mesos leader.

## v1.1.1

- FIX: Increase http client timeout

## v1.1.0

- FIX: mesos cni and docker network alias handling
- FIX: mesos task could be removed after it failed during restart
- DEL: unneeded mesoscni env parameter.
- ADD: Mesos Healthcheck integration.
- CHANGE: optimise offer handling
- ADD: Posibility to add plugins. For an example, take a look into the plugins directory.
- ADD: Kafka Plugin to forware mesos event messages to kafka
- CHANGE: Migrate to google protobuf
- FIX: Scale up/down of a mesos task.
- FIX: unique placement

## 1.0.3

- FIX-SECURITY: https://github.com/AVENTER-UG/mesos-compose/security/dependabot/3

## 1.0.2

- FIX-SECURITY: https://github.com/AVENTER-UG/mesos-compose/security/dependabot/2

## 1.0.1

- FIX-SECURITY: https://github.com/AVENTER-UG/mesos-compose/security/dependabot/1

## 1.0.0

- CHANGE: !!!! Braking changes !!!! Change mesos-compose labels to yaml labels. For an example, please take a look into the docs.
- CHANGE: !!!! Braking changes !!!! Change environment variables format. For an example, please take a look into the docs.
- CHANGE: Move mesos specific functions into own module
- ADD: Implicit reconcile to remove unknown Mesos Tasks
- ADD: Default volume driver parameter `DEFAULT_VOLUME_DRIVER`
- ADD: Set default values for Network (default) and NetworkMode (user)
- CHANGE: set net-alias only if it's defined
- ADD: Support for on-demand scale up of an instance
- ADD: Healthcheck for how much instances are running and if it to less, deploy the missing one
- UPDATE: Optimize revive and suppress
- ADD: Support for docker ulimits memlock and nofile
- FIX: Missing port if we use "UDP" as protocol under ports
- CHANGE: Optimize heartbeat
- ADD: Shell flag to tell Mesos it should treat the command as shell (like: `/bin/sh -c <command> <args>`)
- ADD: API endpoint for reregistration of the framework. These make it possible to force a registration after Mesos lost the framework.
- ADD: Mesos shell flag to control how the containerizer execute the command. If shell is true, then the command will be executes with `/bin/sh -c`.
- FIX: Default CPU ressource
- ADD: API Endpoint to supress the Framework
- ADD: support for wss and h2c
- FIX: Custom Executor command was not set
- ADD: Support Mesos fetch to download files during runtime into the containers sandbox.
- ADD: Scale up and down of mesos tasks
- FIX: Increase TASK_ID during task_lost restart
- FIX: Exit startup if cannot connect redis
- CHANGE: Change the discovery name format to fit the DNS RFC.
- ADD: Better support for mesos containerizer (thanks to @harryzz)
- ADD: Command Attributes (thanks to @harryzz)
- FIX: Conflict between reconcile and heatbeat could end in a task restart loop
- ADD: Parameter to configure the Mesos Task DiscoveryInfoName Delimiter `DISCOVERY_INFONAME_DELIMITER`. Default value is ".".
- ADD: Parameter to configure the Mesos Task DiscoveryPortName Delimiter `DISCOVERY_PORTNAME_DELIMITER`. Default value is "_".
- ADD: Constraint `unique` to run only one instance of a task per node.
- UPDATE: mesos-cli to support the new avmesos-cli python modules.

## 0.4.2

- CHANGE: Command attribute from array to string.
- ADD: Support of docker-compose capability parameter [cap_add](https://docs.docker.com/compose/compose-file/#cap_add).
- ADD: Support of mesos command executor. Set the label:
  `biz.aventer.mesos_compose.container_type: "NONE"`
- ADD: Support of environment variables for executer.
- CHANGE: Optimize offer handling for ports.
- CHANGE: Optimize redis key search.
- ADD: Docker container support for custom executor.
- ADD: Support for docker compose [cap-drop](https://docs.docker.com/compose/compose-file/#cap_drop).
- ADD: Support for docker comport [pull_policy](https://docs.docker.com/compose/compose-file/#pull_policy). Support always (default) and "missing".
- FIX: Recalculate the HostPorts if the Mesos Task is failed.
- ADD: Support for docker compose "placement -> constraints -> node.hostname" command.
- ADD: Resubscription after the connection to mesos master is lost.
- ADD: Mesos CLI Plugin to launch and kill mesos-compose workload.
- ADD: Support for Hashicorp Vault.
- ADD: Overwrite the webui URL by env "FRAMEWORK_WEBUIRUL"
- ADD: Mesos CLI restart and update service.
- ADD: Support for `node.platform.os` and `node.platform.arch` constraint
- ADD: Support of docker-compose command [restart](https://docs.docker.com/compose/compose-file/#read_only)
- ADD: Show all Tasks as API call and mesos-cli command.
- FIX: Offer for multiple host ports.
- FIX: kill services and tasks
- FIX: restart services and tasks to prevent unmanaged tasks
- ADD: Support for user defined network with exposed ports
- ADD: customize taskname `bis.aventer.mesos_compose.taskname: "test:app"`.
- FIX: restore MesosAgent info after update task by API
- FIX: Remove LOST mesos tasks from redis
- ADD: mesos reconcile loop to periodically sync state with mesos

## v0.4.0

- ADD: Redis Connection retry and health check.
- FIX: CountRedisKeys only the own (frameworkName) one.
- CHANGE: Refactoring to be more flexible.

## v0.3.1

- ADD-11: Docker Compose Replicas

## v0.3.1

- FIX-8: Set disk resource to the default value if its unset.
- FIX-7: Set mesos resources at executor info.

## v0.3.0

- Add Redis Authentication Support
- Change DB items framework and framework_config to be saved with the
  frameworkName as prefix.
- The default prefix of hostname and tasks are the FrameworkName.
- Optimize framework suppress
- Add constraint "hostname eq" to accept only offers with the given hostname.
- Check if the random generated port is already in use at the given agent.
- Change TaskID format to be more speakable.
- IMPORTANT!!! Change the API URL's.
- Add force pull image
- Add executor label to define custom mesos executor.

## v0.2.0

- Start with changelog
- Add Docker Ports
- Add Mesos Discovery
- Add Docker Volume support
- Add Environment variables
- Add Restart of Mesos Tasks
- Add Kill Service
- Add Kill single task of a service
- Add Framework Suppress if there is nothing to schedule
- Remove tasks after they are killed
- Add reconcile after frameware resubscribe
- Add TLS Server Support (env variable SSL_CRT_BASE64, SSL_KEY_BASE64)

## v0.1.0

- First inn
