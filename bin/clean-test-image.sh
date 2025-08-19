#!/bin/sh

prefix_image="docker-registry.anandadf.my.id/micros-template/"
service_name="file_service"

echo "Removing test in local" >/dev/stderr
docker rmi "$prefix_image$service_name:test"
