#!/bin/bash

stop() {
    if status ; then
        docker stop $NAME
        echo "$NAME stopped"
    fi
}

status() {
    if docker ps --format "{{.Names}}" | grep ^${NAME}$ >/dev/null ; then
        echo "$NAME running"
    else
        echo "$NAME not running"
        return 1
    fi
    return 0
}

docker_ip() {
    docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $NAME    
}
