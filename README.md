# Sample controller

A super basic kubernetes controller that annotates the creation time of a pod. This controller also has leader election
working via setting the correct options in the manager options struct.

## building

run `make build` to compile the application (the binary name is `pod-timestamp-controller`)

## testing

TODO

## Running locally

To run the controller locally you need [kind](https://github.com/kubernetes-sigs/kind) installed (it requires you to
have docker installed).

1. start up the cluster and registry by running `make create-local-env`

2. build the image by running `make docker-build TAG=<version number>`

3. push the image by running `make docker-push-local TAG=<version number>`

5. deploy the helm chart by running `make install TAG=<version number>`

6. view the logs of the running controller by `kubectl logs -n sample-controller <controller pod name>`

## Clean up

1. run `make destroy-local-env` to destroy the cluster and local registry