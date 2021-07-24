package reconcile

import (
	"context"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Loop struct{
	Client client.Client
	Log logr.Logger
}

func (l *Loop) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := l.Log.WithValues("pod", request.NamespacedName)
	var pod v1.Pod

	getPodErr := l.Client.Get(ctx, request.NamespacedName, &pod)
	if getPodErr != nil {
		// TODO remove not found pods from queue
		// error package here (need to bookmark this) https://pkg.go.dev/k8s.io/apimachinery/pkg/api/errors
		log.Error(getPodErr, "could not get pod")
		return reconcile.Result{}, getPodErr
	}
	log.Info("pod is here", "pod name", request.Name, "timestamp of creation", pod.Status.StartTime.Time.String())
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["creation"] = pod.Status.StartTime.Time.String()
	updateErr := l.Client.Update(ctx, &pod)
	if updateErr != nil {
		log.Error(updateErr, "could not update timestamp")
		return reconcile.Result{}, updateErr
	}
	return reconcile.Result{}, nil
}
