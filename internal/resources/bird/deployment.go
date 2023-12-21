package bird

import (
	"context"
	"fmt"

	"praios.lf-net.org/littlefox/router-operator/internal/resources"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func reconcileDeployments(ctx context.Context, c client.Client, reconciliation *resources.RouterReconciliation) error {
	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: reconciliation.Router.GetNamespace(),
			Name:      fmt.Sprintf("router-%s-bird", reconciliation.Router.GetName()),
		},
	}

	status, err := controllerutil.CreateOrUpdate(ctx, c, &deploy, func() error {
		deploy.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			*metav1.NewControllerRef(reconciliation.Router, reconciliation.Router.GetObjectKind().GroupVersionKind()),
		}

		for label, value := range reconciliation.Router.ObjectMeta.GetLabels() {
			metav1.SetMetaDataLabel(&deploy.ObjectMeta, label, value)
		}

		for annotation, value := range reconciliation.Router.ObjectMeta.GetAnnotations() {
			metav1.SetMetaDataAnnotation(&deploy.ObjectMeta, annotation, value)
		}

		selector := metav1.LabelSelector{}
		if err := metav1.Convert_Map_string_To_string_To_v1_LabelSelector(&deploy.ObjectMeta.Labels, &selector, nil); err != nil {
			return fmt.Errorf("error converting labels to selector: %w", err)
		}

		deploy.Spec = appsv1.DeploymentSpec{
			Selector: &selector,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      deploy.Labels,
					Annotations: deploy.Annotations,
				},
				Spec: corev1.PodSpec{
					NodeSelector: reconciliation.Router.Spec.NodeSelector,
					Tolerations:  reconciliation.Router.Spec.Tolerations,
					HostNetwork:  true,
					Containers: []corev1.Container{
						{
							Name:  "bird",
							Image: "vnxme/bird:2.13",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/bird/",
									ReadOnly:  true,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"CAP_NET_ADMIN",
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("router-%s-config", reconciliation.Router.Name),
									},
								},
							},
						},
					},
				},
			},
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating or updating Deployment: %w", err)
	}

	log.FromContext(ctx).Info("Deployment reconciled", "deployment", client.ObjectKeyFromObject(&deploy), "status", status)
	return nil
}
