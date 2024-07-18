#!/usr/bin/env python3

import sys
import json
import docker


def parse_payload():
    payload = sys.argv[1]
    payload_parsed = json.loads(payload)
    version = payload_parsed["version"]
    env = payload_parsed["env"]
    print(f"Getting {version}...")
    return version, env


def main():
    version, env = parse_payload()

    if version == "":
        print("Version is empty, getting latest version...")
        version = "latest"

    if env == "prod":
        print("Redeploying production container...")
        name = "elrek-system_go_prod"
        volumes = {"logs": {"bind": "/app/logs", "mode": "rw"}}
        environment = [
            "DB_USERNAME=<DB_USERNAME>",
            "DB_PASSWORD=<DB_PASSWORD>",
            "DB_IP=<DB_IP>",
            "DB_PORT=<DB_PORT>",
            "DB_NAME=<DB_NAME>",
        ]
        ports = {"3000/tcp": "<PORT_AS_INT>"}
    else:
        print("Redeploying development container...")
        name = "elrek-system_go_dev"
        volumes = {"logs": {"bind": "/app/logs", "mode": "rw"}}
        environment = [
            "DB_USERNAME=<DB_USERNAME>",
            "DB_PASSWORD=<DB_PASSWORD>",
            "DB_IP=<DB_IP>",
            "DB_PORT=<DB_PORT>",
            "DB_NAME=<DB_NAME>",
        ]
        ports = {"3000/tcp": "<PORT_AS_INT>"}

    client = docker.from_env()
    try:
        container = client.containers.get(name)
        container.stop()
        print("Stopped current version")
    except:
        pass
    client.containers.prune()
    print("Removed current version")

    client.containers.run(
        f"sc4n1a471/elrek-system_go:{version}",
        detach=True,
        volumes=volumes,
        environment=environment,
        ports=ports,
        name=name,
        restart_policy={"Name": "on-failure", "MaximumRetryCount": 5},
    )
    print(f"Version {version} was deployed successfully")


if __name__ == "__main__":
    main()
