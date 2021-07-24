package main

import (
	"os"

	"github.com/ladyserena/sample-controller/prestart"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {

	mgr := prestart.Do()

	if startErr := mgr.Start(signals.SetupSignalHandler()); startErr != nil {
		mgr.GetLogger().Error(startErr, "could not keep running manager")
		os.Exit(1)
	}

}
