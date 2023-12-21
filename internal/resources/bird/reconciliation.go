package bird

import (
	"context"
	"fmt"

	"praios.lf-net.org/littlefox/router-operator/internal/resources"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ReconcileAll(ctx context.Context, c client.Client, reconciliation *resources.RouterReconciliation) error {
	reconcilers := map[string]func(context.Context, client.Client, *resources.RouterReconciliation) error{
		"ConfigMaps":  reconcileConfigMaps,
		"Deployments": reconcileDeployments,
	}

	for name, reconciler := range reconcilers {
		if err := reconciler(ctx, c, reconciliation); err != nil {
			return fmt.Errorf("error reconciling %v: %v", name, err)
		}
	}

	return nil
}
