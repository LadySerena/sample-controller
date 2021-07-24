package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/ladyserena/sample-controller/reconcile"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func Bail(logger logr.Logger, err error, message string) {
	logger.Error(err, message)
	os.Exit(1)
}

func main() {
	log := zap.New().WithName("sample-controller")
	logf.SetLogger(log)
	restConfig, configErr := config.GetConfig()
	if configErr != nil {
		errorMessage := "error creating restConfig from pod service account"
		if configErr == rest.ErrNotInCluster {
			errorMessage = "this needs to run within the cluster"
		}
		Bail(log, configErr, errorMessage)
	}

	controllerScheme := runtime.NewScheme()
	schemeErr := corev1.AddToScheme(controllerScheme)
	if schemeErr != nil {
		Bail(log, schemeErr, "could not add to scheme")
	}

	managerOptions := manager.Options{
		Scheme:                        controllerScheme,
		MapperProvider:                nil,
		SyncPeriod:                    nil,
		Logger:                        log,
		LeaderElection:                false,
		LeaderElectionResourceLock:    "",
		LeaderElectionNamespace:       "",
		LeaderElectionID:              "",
		LeaderElectionConfig:          nil,
		LeaderElectionReleaseOnCancel: false,
		LeaseDuration:                 nil,
		RenewDeadline:                 nil,
		RetryPeriod:                   nil,
		Namespace:                     "",
		MetricsBindAddress:            "0.0.0.0:8080", // todo make configurable
		HealthProbeBindAddress:        "0.0.0.0:8081", // todo make configurable
		ReadinessEndpointName:         "",
		LivenessEndpointName:          "",
		Port:                          0,
		Host:                          "",
		CertDir:                       "",
		WebhookServer:                 nil,
		NewCache:                      nil,
		NewClient:                     nil,
		ClientDisableCacheFor:         nil,
		DryRunClient:                  false,
		EventBroadcaster:              nil,
		GracefulShutdownTimeout:       nil,
		Controller:                    v1alpha1.ControllerConfigurationSpec{},
	}

	mgr, managerCreateError := manager.New(restConfig, managerOptions)
	if managerCreateError != nil {
		Bail(log, managerCreateError, "could not create the manager")
	}

	mgr.AddHealthzCheck("test", func(req *http.Request) error {
		return nil
	})
	mgr.AddReadyzCheck("test2", func(req *http.Request) error {
		return nil
	})
	podController, controllerCreateError := controller.New("sample-controller", mgr, controller.Options{
		MaxConcurrentReconciles: 1,
		Reconciler: &reconcile.Loop{
			Client: mgr.GetClient(),
			Log: log,
		},
		RateLimiter:             nil,
		Log:                     log,
		CacheSyncTimeout:        2 * time.Minute, // todo make configurable
	})
	if controllerCreateError != nil {
		Bail(log, controllerCreateError, "could not create controller")
	}

	watchErr := podController.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
	if watchErr != nil {
		Bail(log, watchErr, "unable to watch pods")
	}

	if startErr := mgr.Start(signals.SetupSignalHandler()); startErr != nil {
		Bail(log, startErr, "could not keep running manager")
	}

}
