package init

import (
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/ladyserena/sample-controller/reconcile"
	"go.uber.org/zap/zapcore"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	k8szap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Do runs the init process of all the components for the controller to run
func Do() manager.Manager {
	logger := setGlobalLogger()
	restConfig, configErr := config.GetConfig()
	if configErr != nil {
		bail(configErr, "could not create rest restConfig")
	}
	scheme, schemeErr := getScheme()
	if schemeErr != nil {
		bail(schemeErr, "could not create scheme")
	}
	mgr, managerCreateErr := createManagerFromEnvironment(scheme, restConfig)
	if managerCreateErr != nil {
		bail(managerCreateErr, "could not create mgr")
	}

	ctrl, ctrlErr := createController(logger, mgr)
	if ctrlErr != nil {
		bail(ctrlErr, "could not create controller")
	}

	watchErr := ctrl.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
	if watchErr != nil {
		bail(watchErr, "could not watch pods")
	}

	return mgr
}

func setGlobalLogger() logr.Logger {
	timestampEncoder := k8szap.EncoderConfigOption(func(encoderConfig *zapcore.EncoderConfig) {
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	})

	log := k8szap.New(k8szap.WriteTo(os.Stdout), k8szap.JSONEncoder(timestampEncoder)).WithName("sample-controller")
	logf.SetLogger(log)
	return log
}

func getScheme() (*runtime.Scheme, error) {
	controllerScheme := runtime.NewScheme()
	schemeErr := corev1.AddToScheme(controllerScheme)
	return controllerScheme, schemeErr
}

func createManagerFromEnvironment(scheme *runtime.Scheme, config *rest.Config) (manager.Manager, error) {
	metricsAddress := os.Getenv("METRICS_BIND_ADDRESS")
	if metricsAddress == "" {
		metricsAddress = ":8080"
	}
	healthAddress := os.Getenv("HEALTH_BIND_ADDRESS")
	if healthAddress == "" {
		healthAddress = ":8081"
	}

	options := manager.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddress,
		HealthProbeBindAddress: healthAddress,
	}

	mgr, managerErr := manager.New(config, options)
	if managerErr != nil {
		return nil, managerErr
	}

	addHealthErr := mgr.AddHealthzCheck("health", healthCheckFunc)
	if addHealthErr != nil {
		return nil, addHealthErr
	}

	addReadyErr := mgr.AddReadyzCheck("ready", healthCheckFunc)
	if addReadyErr != nil {
		return nil, addReadyErr
	}

	return mgr, nil
}

func createController(logger logr.Logger, mgr manager.Manager) (controller.Controller, error) {
	return controller.New("sample-controller", mgr, controller.Options{
		MaxConcurrentReconciles: 0,
		Reconciler: &reconcile.Loop{
			Client: mgr.GetClient(),
			Log:    logger,
		},
		Log: logger,
	})
}

func bail(err error, message string) {
	logf.Log.Error(err, message)
	os.Exit(1)
}

func healthCheckFunc(req *http.Request) error {
	return nil
}
