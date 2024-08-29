"""
The mesos-compose plugin.
"""

import toml
import json

import urllib3

from avmesos import cli
from avmesos.cli.exceptions import CLIException
from avmesos.cli.plugins import PluginBase
from avmesos.cli.util import Table
from avmesos.cli.mesos import get_frameworks, get_framework_address
from avmesos.cli import http


PLUGIN_NAME = "compose"
PLUGIN_CLASS = "mesosCompose"

VERSION = "0.1.0"

SHORT_HELP = "Interacts with the Mesos-Compose Framework"

class Config():

    def __init__(self, main):
        """
        Get authentication header for the framework
        """

        self.main = main

        try:
            data = toml.load(self.main.config.path)
        except Exception as exception:
            raise CLIException(
                "Error loading config file as TOML: {error}".format(
                    error=exception)
            ) from exception

        self.data = data["compose"].get(self.main.framework_name)

    def principal(self):
        """
        Return the principal in the configuration file
        """
        return self.data.get("principal")

    def secret(self):
        """
        Return the secret in the configuration file
        """

        return self.data.get("secret")

    def ssl_verify(self, default=False):
        """
        Return if the ssl certificate should be verified
        """
        ssl_verify = self.data.get("ssl_verify", default)
        if not isinstance(ssl_verify, bool):
            raise CLIException("The 'ssl_verify' field must be True/False")

        return ssl_verify

    # pylint: disable=no-self-use
    def agent_timeout(self, default=5):
        """
        Return the connection timeout of the agent
        """

        return default



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
            "flags": {},
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
            "arguments": ["<framework-name>", "<task>"],
            "flags": {},
            "short_help": "Kill a single task (ID) or a whole service (Task Name)",
            "long_help": "Use the \"ID\" to Kill a single task or the \"Task Name\" to kill the entire service."
        },
        "restart": {
            "arguments": ["<framework-name>", "<task>"],
            "flags": {},
            "short_help": "Restart a single task (ID) or a whole service (Task Name)",
            "long_help": "Use the \"ID\" to restart a single task or the \"Task Name\" to restart the entire service."
        },
        "framework": {
            "arguments": ["<framework-name>", "<operations>"],
            "flags": {},
            "short_help": "Framework Commands.",
            "long_help": "Framework Commands\n\treregister - force the reregistration of the framework. !!! ONLY USE IT DURING MESOS CONNECTION ERRORS. !!!\n\tsuppress - supress the framework\n",
        }
    }

    def framework(self, argv):
        """
        Framework commands.
        """

        try:
            master = self.config.master()
            config = self.onfig
            # pylint: disable=attribute-defined-outside-init
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        if argv["<operations>"] == "reregister":
            print("Force the reregistration of these framework")

            framework_address = get_framework_address(
                self.get_framework_id(argv), master, config
            )
            data = self.write_endpoint(
                framework_address,
                "/api/compose/v0/framework/reregister",
                self,
                "PUT"
            )
            print(data)

        if argv["<operations>"] == "suppress":
            print("Supress the Framework")

            framework_address = get_framework_address(
                self.get_framework_id(argv), master, config
            )
            data = self.write_endpoint(
                framework_address,
                "/api/compose/v0/framework/supress",
                self,
                "PUT"
            )
            print(data)


    def launch(self, argv):
        """
        Launch Mesos workload from compose file
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
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
            message = json.loads(
                self.write_endpoint(
                    framework_address,
                    "/api/compose/v0/" + project,
                    self.mesos_config,
                    "PUT",
                    filename,
                )
            )

            try:
                print(json.dumps(message, indent=2, ensure_ascii=False))
            except Exception as exception:
                print(message)
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
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        self.framework_name = argv["<framework-name>"]

        if self.framework_name is not None:

            framework_address = get_framework_address(
                self.get_framework_id(argv), master, config
            )

            message = json.loads(
                http.read_endpoint(framework_address, "/api/compose/v0/tasks", self.mesos_config)
            )

            try:
                if not message:
                    print("There are no tasks running in the cluster.")
                    return

                table = Table(["ID", "Task Name", "State", "Mesos Agent"])
                for task in message:

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
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
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
                    self.mesos_config,
                    "UPDATE",
                    filename,
                )
            )

            try:
                print(json.dumps(data, indent=2, ensure_ascii=False))
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
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        framework_address = get_framework_address(
            self.get_framework_id(argv), master, config
        )

        data = http.read_endpoint(framework_address, "/api/compose/versions", self.mesos_config)

        print(data)

    def kill(self, argv):
        """
        Kill mesos task
        """

        try:
            master = self.config.master()
            config = self.config
            # pylint: disable=attribute-defined-outside-init
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        framework_address = get_framework_address(
            self.get_framework_id(argv), master, config
        )

        if argv.get("<task>").__contains__(":"):
            task = argv.get("<task>").split(":")
            project = task[1]
            service = task[2]
            data = self.write_endpoint(
                framework_address,
                "/api/compose/v0/" + project + "/" + service,
                self.mesos_config,
                "DELETE",
            )
        else:
            data = self.write_endpoint(
                framework_address,
                "/api/compose/v0/tasks/"+ argv.get("<task>"),
                self.mesos_config,
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
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
        except Exception as exception:
            raise CLIException(
                "Unable to get leading master address: {error}".format(error=exception)
            ) from exception

        framework_address = get_framework_address(
            self.get_framework_id(argv), master, config
        )

        if argv.get("<task>").__contains__(":"):
            """
            The given parameter looks like a task-name
            """
            task = argv.get("<task>").split(":")
            project = task[1]
            service = task[2]

            data = self.write_endpoint(
                framework_address,
                "/api/compose/v0/" + project + "/" + service + "/restart",
                self.mesos_config,
                "PUT",
            )
        else:
            data = self.write_endpoint(
                framework_address,
                "/api/compose/v0/tasks/"+ argv.get("<task>") + "/restart",
                self.mesos_config,
                "PUT",
            )

        print(data)

    def info(self, argv):
        """
        Get information about the running Mesos compose framework
        """

        try:
            master = self.config.master()
            # pylint: disable=attribute-defined-outside-init
            self.framework_name = argv["<framework-name>"]
            self.mesos_config = Config(self)
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
            if config.ssl_verify() is not True:
                http = urllib3.PoolManager(cert_reqs='CERT_NONE')
            else:
                http = urllib3.PoolManager()
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
