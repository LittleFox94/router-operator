package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CommonCondition struct {
	// Status of this condition.
	Status corev1.ConditionStatus `json:"status,omitempty"`

	// Reason for this condition being this state, computer-readable CamelCaseString.
	Reason string `json:"reason,omitempty"`

	// Human-readable message to explain why this condition is in this state.
	Message string `json:"message,omitempty"`

	// Timestamp of this condition last changing its status.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// Timestamp of this condition being last updated to actual state.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}
