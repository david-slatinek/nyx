#!/usr/bin/env bash

while getopts a:d: flag; do
  case "${flag}" in
  a) aws_url=${OPTARG} ;;
  d) docker_url=${OPTARG} ;;
  *) echo "Invalid flag" && exit 1 ;;
  esac
done

if [ -z "$aws_url" ] || [ -z "$docker_url" ]; then
  echo "Invalid arguments, must be -a aws_url AND -d docker_url"
  exit 1
fi

aws_services=("${aws_url}/couchdb" "${aws_url}/dialog-api" "${aws_url}/website" "${aws_url}/summary-api"
  "${aws_url}/categories-api" "${aws_url}/recommend-db" "${aws_url}/recommend-api" "${aws_url}/monitoring")

for service in "${aws_services[@]}"; do
  docker pull "$service"
  echo -e "\tPulled $service\n"
done

docker_services=("${docker_url}/dialog-model" "${docker_url}/summary-model" "${docker_url}/recommend-model")

for service in "${docker_services[@]}"; do
  docker pull "$service"
  echo -e "\tPulled $service\n"
done
