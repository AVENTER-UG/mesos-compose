# Changelog

## Master

- Add Redis Authentication Support
- Change DB items framework and framework_config to be saved with the
  frameworkName as prefix.
- The default prefix of hostname and tasks are the FrameworkName.
- Optimize framework suppress
- Add constraint "hostname eq" to accept only offers with the given hostname.
- Check if the random generated port is already in use at the given agent.
- IMPORTANT!!! Change the API URL's.

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
