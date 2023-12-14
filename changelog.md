# Changelog

## master

- FIX: mesos cni and docker network alias handling
- FIX: mesos task could be removed after it failed during restart
- DEL: unneeded mesoscni env parameter.
- ADD: Mesos Healthcheck integration. 

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
