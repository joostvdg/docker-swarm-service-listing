#!/usr/bin/env bash
docker build -t flow-proxy-service-lister .
docker run --rm  -p 7777:7777 -v /var/run/docker.sock:/var/run/docker.sock flow-proxy-service-lister