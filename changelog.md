# Changelog

## master

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
- ADD: Support for Hashicorp Vault also for the environment variables of the framework.
- ADD: Mesos CLI restart and update service.
- ADD: Support for `node.platform.os` and `node.platform.arch` constraint
- ADD: Support of docker-compose command [restart](https://docs.docker.com/compose/compose-file/#read_only) 
  
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
