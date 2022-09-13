"""
The mesos-compose plugin.
"""

import toml
import cli
import json

from urllib.parse import urlencode
import urllib3

from cli.exceptions import CLIException
from cli.plugins import PluginBase
from cli.util import Table

from cli.mesos import get_frameworks, get_framework_address
from cli import http


PLUGIN_NAME = "compose"
PLUGIN_CLASS = "mesosCompose"

VERSION = "0.1.0"

SHORT_HELP = "Interacts with the Mesos-Compose Framework"


class mesosCompose(PluginBase):
    """
    The mesos-compose plugin.
    """

    COMMANDS = {
        "version": {
            "arguments": ["<framework-name>"],
            "flags": {},
            "short_help": "Get the version number of Mesos compose",
            "long_help": "Get the version number of Mesos compose",
        },
        "info": {
            "arguments": ["<framework-name>"],
            "flags": {},
            "short_help": "Get information about the running Mesos compose framework.",
            "long_help": "Get information about the running Mesos compose framework.",
        },
        "list": {
            "arguments": ["<framework-name>"],
            "flags": {"-a --all": "list all tasks, not only running [default: False]"},
            "short_help": "Show all running tasks.",
            "long_help": "Show all running tasks.",
        },
        "launch": {
            "arguments": ["<framework-name>", "<project>", "<compose-file>"],
            "flags": {},
            "short_help": "Launch Mesos workload from compose file",
            "long_help": "Launch Mesos workload from compose file",
        },
        "update": {
            "arguments": ["<framework-name>", "<project>", "<compose-file>"],
            "flags": {},
            "short_help": "Update service from compose file",
            "long_help": "Update service from compose file",
        },
        "kill": {
            "arguments": ["<framework-name>", "<project>", "<service>"],
            "flags": {},
            "short_help": "Kill Mesos compose workload",
            "long_help": "Kill Mesos compose workload",
        },
        "restart": {
            "arguments": ["<framework-name>", "<project>", "<service>"],
            "flags": {},
            "short_help": "Restart service",
            "long_help": "Restart service",
        },
    }

    def launch(self, argv):
        """
        Launch Mesos workload from compose file
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        project = argv["<project>"]
        filename = argv["<compose-file>"]
        self.framework_name = argv["<framework-name>"]

        if (
            project is not None
            and filename is not None
            and self.framework_name is not None
        ):
            print("Launch workload " + project)

            framework_address = get_framework_address(
                self.get_framework_id(argv), master, config
            )
            data = json.loads(
                self.write_endpoint(
                    framework_address,
                    "/api/compose/v0/" + project,
                    self,
                    "PUT",
                    filename,
                )
            )

            try:
                message = json.loads(data["Message"])
                print(json.dumps(message, indent=2, ensure_ascii=False))
            except Exception as exception:
                print(data)
        else:
            print("Nothing to Launch")

    def list(self, argv):
        """
        Show running tasks
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        self.framework_name = argv["<framework-name>"]

        if self.framework_name is not None:

            framework_address = get_framework_address(
                self.get_framework_id(argv), master, config
            )

            data = json.loads(
                http.read_endpoint(framework_address, "/api/compose/v0/tasks", self)
            )

            try:
                message = json.loads(data["Message"])

                if not message:
                    print("There are no tasks running in the cluster.")
                    return

                table = Table(["ID", "Task Name", "State", "Mesos Agent"])
                for task in message:

                    if not argv["--all"] and task["State"] != "TASK_RUNNING":
                        continue

                    table.add_row(
                        [
                            task["TaskID"],
                            task["task_name"],
                            task["State"],
                            task["MesosAgent"].get("hostname"),
                        ]
                    )

            except Exception as exception:
                raise CLIException(
                    "Unable to build table of tasks: {error}".format(error=exception)
                )

            print(str(table))

    def update(self, argv):
        """
        Update running service from compose file
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        project = argv["<project>"]
        filename = argv["<compose-file>"]
        self.framework_name = argv["<framework-name>"]

        if (
            project is not None
            and filename is not None
            and self.framework_name is not None
        ):
            print("Update workload " + project)

            framework_address = get_framework_address(
                self.get_framework_id(argv), master, config
            )
            data = json.loads(
                self.write_endpoint(
                    framework_address,
                    "/api/compose/v0/" + project,
                    self,
                    "UPDATE",
                    filename,
                )
            )

            try:
                message = json.loads(data["Message"])
                print(json.dumps(message, indent=2, ensure_ascii=False))
            except Exception as exception:
                print(data)

        else:
            print("Nothing to Update")

    def version(self, argv):
        """
        Get the version information of Kubernetes
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
            self.framework_name = argv["<framework-name>"]
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        framework_address = get_framework_address(
            self.get_framework_id(argv), master, config
        )

        data = http.read_endpoint(framework_address, "/api/compose/versions", self)

        print(data)

    def kill(self, argv):
        """
        Kill mesos task
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
            self.framework_name = argv["<framework-name>"]
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        project = argv.get("<project>")
        service = argv.get("<service>")
        framework_address = get_framework_address(
            self.get_framework_id(argv), master, config
        )

        data = self.write_endpoint(
            framework_address,
            "/api/compose/v0/" + project + "/" + service,
            self,
            "DELETE",
        )

        print(data)

    def restart(self, argv):
        """
        Restart mesos task
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
            self.framework_name = argv["<framework-name>"]
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        project = argv.get("<project>")
        service = argv.get("<service>")
        framework_address = get_framework_address(
            self.get_framework_id(argv), master, config
        )

        data = self.write_endpoint(
            framework_address,
            "/api/compose/v0/" + project + "/" + service + "/restart",
            self,
            "PUT",
        )

        print(data)

    def info(self, argv):
        """
        Get information about the running Mesos compose framework
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.mesos_config = self._get_config()
            self.framework_name = argv["<framework-name>"]
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        framework_address = get_framework_address(
            self.get_framework_id(argv), master, self
        )
        print("Framework Address:               " + framework_address)
        print("Framework ID:                    " + self.get_framework_id(argv))
        print("\n")

    def get_framework_id(self, argv):
        """
        Resolv the id of a framework by the name of a framework
        """

        """Check if the given parameter is a ID or a name"""
        if argv["<framework-name>"].count("-") != 5:
            data = get_frameworks(self.config.master(), self.config)
            for framework in data:
                if (
                    framework["active"] is not True
                    or framework["name"].lower() != argv["<framework-name>"].lower()
                ):
                    continue
                return framework["id"]
        return argv["<framework-name>"]

    def principal(self):
        """
        Return the principal in the configuration file
        """

        return self.mesos_config["compose"].get(self.framework_name).get("principal")

    def secret(self):
        """
        Return the secret in the configuration file
        """

        return self.mesos_config["compose"].get(self.framework_name).get("secret")

    # pylint: disable=no-self-use
    def agent_timeout(self, default=5):
        """
        Return the connection timeout of the agent
        """

        return default

    def _get_config(self):
        """
        Get authentication header for the framework
        """

        try:
            data = toml.load(self.config.path)
        except Exception as exception:
            raise CLIException(
                "Error loading config file as TOML: {error}".format(error=exception)
            ) from exception

        return data

    def write_endpoint(self, addr, endpoint, config, method, filename=None):
        """
        Read the specified endpoint and return the results.
        """

        try:
            addr = cli.util.sanitize_address(addr)
        except Exception as exception:
            raise CLIException(
                "Unable to sanitize address '{addr}': {error}".format(
                    addr=addr, error=str(exception)
                )
            )
        try:
            url = "{addr}{endpoint}".format(addr=addr, endpoint=endpoint)
            if config.principal() is not None and config.secret() is not None:
                headers = urllib3.make_headers(
                    basic_auth=config.principal() + ":" + config.secret(),
                )
            else:
                headers = None
            http = urllib3.PoolManager()
            content = ""
            if filename is not None:
                data = open(filename, "rb")
                content = data.read()
            http_response = http.request(
                method,
                url,
                headers=headers,
                body=content,
                timeout=config.agent_timeout(),
            )
            return http_response.data.decode("utf-8")

        except Exception as exception:
            raise CLIException(
                "Unable to open url '{url}': {error}".format(
                    url=url, error=str(exception)
                )
            )
