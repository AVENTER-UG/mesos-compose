# Changelog

## master

- CHANGE: Command attribute from array to string. 
- ADD: Support of docker-compose capability parameter "cap_add".
- ADD: Support of mesos command executor. Set the label:
  `biz.aventer.mesos_compose.container_type: "NONE"`
- ADD: Support of environment variables for executier.
  
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
