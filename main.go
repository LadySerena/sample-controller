package main

import (
	"os"

	"github.com/ladyserena/sample-controller/init"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {

	mgr := init.Do()

	if startErr := mgr.Start(signals.SetupSignalHandler()); startErr != nil {
		mgr.GetLogger().Error(startErr, "could not keep running manager")
		os.Exit(1)
	}

}
