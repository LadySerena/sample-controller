# Sample controller

A super basic kubernetes controller that annotates the creation time of a pod.

## building

run `make build` to compile the application (the binary name is `pod-timestamp-controller`)

## testing

TODO

## Running locally

To run the controller locally you need [kind](https://github.com/kubernetes-sigs/kind) installed (it requires you to
have docker installed).

1. start up the cluster and registry by running `./scripts/create-kind-with-registry.sh`

2. build the image by running `docker build -t localhost:5000/sample-controller:<version number> .`

3. push the image by running `docker push localhost:5000/sample-controller:<version number>`

4. check that the `values.yaml` file ([here](./config/sample-controller/values.yaml)) has the matching tag for the image
   tag

5. deploy the helm chart by running `helm install sample-controller ./config/sample-controller`

6. view the logs of the running controller by `kubectl logs -n default <controller pod name>`

## Clean up

1. run `./scripts/destroy-kind-and-registry.sh` to destroy the cluster and local registry