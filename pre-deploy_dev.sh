#!/bin/bash
if [ $(docker ps -a -q -f name=elrek-system_go_dev) ]; then
    docker rm -f elrek-system_go_dev
    echo "Container elrek-system_go_dev deleted."
else
    echo "Container elrek-system_go_dev does not exist."
fi

if [ $(docker images -q sc4n1a471/elrek-system_go:$version-dev) ]; then
    docker rmi -f sc4n1a471/elrek-system_go:$version-dev
    echo "Image sc4n1a471/elrek-system_go:$version-dev deleted."
else
    echo "Image sc4n1a471/elrek-system_go:$version-dev does not exist."
fi