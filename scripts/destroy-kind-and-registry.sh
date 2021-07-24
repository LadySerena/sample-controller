#!/usr/bin/env bash

kind delete cluster --name kind
docker stop kind-registry
docker rm kind-registry