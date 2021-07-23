package main

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	//https://pkg.go.dev/k8s.io/client-go@v0.21.3/rest#InClusterConfig rest config within cluster as pod (gonna yeet rbac into the helm chart)
	mgr := manager.New() // need to sort out the rest config
	_, err := controller.New("sample-controller")

}
