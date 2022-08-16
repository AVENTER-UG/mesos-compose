# Mesos-Compose CLI Installation for Mesos-CLI

If you do not already have installe the mesos cli, please follow the steps under "Install Mesos-CLI" first.

The installation of the Mesos-Compose plugin for mesos-cli is done in few steps.

First, edit the mesos-cli config file.

```bash

vim .mesos/config.toml

```

Add the absolute path of the plugin into the plugin array:

```bash

# The `plugins` array lists the absolute paths of the
# plugins you want to add to the CLI.
plugins = [
  "/example/mesos-compose/mesos_cli/compose"
]

[compose-<FRAMEWORK_PREFIX>]
  principal = "<framework username>"
  secret = "<framework password>"

```

A example for multiple mesos-compose frameworks:

```bash

[compose-mc-a]
  principal = "<framework username>"
  secret = "<framework password>"

[compose-mc-b]
  principal = "<framework username>"
  secret = "<framework password>"

```

As you can see, the "compose" section have to extend with the prefix of Mesos-Compose framework.

Now we will see the M3s plugin in mesos cli:

```bash

mesos-cli help
Mesos CLI

Usage:
  mesos (-h | --help)
  mesos --version
  mesos <command> [<args>...]

Options:
  -h --help  Show this screen.
  --version  Show version info.

Commands:
  agent      Interacts with the Mesos agents
  compose    Interacts with the Mesos-Compose Framework
  config     Interacts with the Mesos CLI configuration file
  framework  Interacts with the Mesos Frameworks
  m3s        Interacts with the Kubernetes Framework M3s
  task       Interacts with the tasks running in a Mesos cluster

```

## Install Mesos-CLI

Download the mesos-cli binary for linux from [here](https://www.aventer.biz/files/sw/Linux/mesos-cli.zip). Extract 
the mesos-cli and copy the file into your PATH directory.
