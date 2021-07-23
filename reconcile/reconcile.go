package reconcile

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Loop struct{}

func (Loop) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	panic("implement me")
}
