package reconcile

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	timestampKey = "creation"
	selectionKey = "timestamp"
)

// Loop is the struct that implements the reconcile.Reconciler interface in controller runtime
type Loop struct {
	Client client.Client
	Log    logr.Logger
}

// Reconcile is where the logic for the reconcilliation loop happens. In this case it is annotating and logging timestamps of pod creation
func (l *Loop) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := l.Log.WithValues("pod", request.NamespacedName)
	var pod v1.Pod

	getPodErr := l.Client.Get(ctx, request.NamespacedName, &pod)
	if getPodErr != nil {
		if errors.IsNotFound(getPodErr) {
			return reconcile.Result{}, nil // don't return error since it's probably it's cache sync issue
		}
		log.Error(getPodErr, "could not get pod")
		return reconcile.Result{}, getPodErr
	}
	shouldAnnotate, checkErr := annotatePod(pod) // check if we should annotate the pod
	if checkErr != nil {
		log.Error(checkErr, "could not check annotations of pod skipping", "pod name", request.Name)
		return reconcile.Result{}, checkErr
	}
	if shouldAnnotate {
		// if we should annotate then log the pod name and the creation timestamp
		log.Info("pod is here", "pod name", request.Name, "timestamp of creation", pod.Status.StartTime.Time.String())
		// only add the timestamp if it isn't already there
		_, keyPresent := pod.Annotations[timestampKey]
		if !keyPresent {
			pod.Annotations[timestampKey] = pod.Status.StartTime.Time.String()
			updateErr := l.Client.Update(ctx, &pod)
			if updateErr != nil {
				log.Error(updateErr, "could not update timestamp")
				return reconcile.Result{}, updateErr
			}
		}
	}
	return reconcile.Result{}, nil
}

func annotatePod(pod v1.Pod) (bool, error) {
	if pod.Annotations == nil {
		return false, nil
	}
	value, isPresent := pod.Annotations[selectionKey]
	if !isPresent {
		return false, nil
	}
	return strconv.ParseBool(value)
}
