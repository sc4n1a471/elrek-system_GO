#!/bin/bash
if [ $(docker ps -a -q -f name=elrek-system_go_prod) ]; then
    docker rm -f elrek-system_go_prod
    echo "Container elrek-system_go_prod deleted."
else
    echo "Container elrek-system_go_prod does not exist."
fi

if [ $(docker images -q sc4n1a471/elrek-system_go:$version) ]; then
    docker rmi -f sc4n1a471/elrek-system_go:$version"
    echo "Image sc4n1a471/elrek-system_go:$version deleted."
else
    echo "Image sc4n1a471/elrek-system_go:$version does not exist."
fi