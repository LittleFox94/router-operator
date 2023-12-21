/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kstatus "sigs.k8s.io/cli-utils/pkg/kstatus/status"
)

// RouterSpec defines the desired state of a Router instance
type RouterSpec struct {
	// NodeSelector is set on the pod create for this Router to pin it to a specific node.
	NodeSelector map[string]string `json:"nodeSelector"`

	Tolerations []corev1.Toleration `json:"tolerations"`

	// World-wide unique ID of this router, usually one of its IPv4 addresses.
	NodeID RouterNodeID `json:"nodeID"`
}

// RouterNodeID defines the world-wide unique ID of a router.
type RouterNodeID struct {
	ID string `json:"id"`
}

// RouterStatus defines the observed state of Router
type RouterStatus struct {
	// Generation of this Router at the time of updating this status.
	ObservedGeneration int64 `json:"observedGeneration"`

	// ClusterAnnouncements is the last number of Announcement objects observed for this Router.
	ClusterAnnouncements uint `json:"clusterAnnouncements"`

	// ClusterSessions is the last number of Session objects observed for this Router.
	ClusterSessions uint `json:"clusterSessions"`

	// ReadyAnnouncements is the last number of Announcement objects that are ready (= exported) on the actual routing daemon.
	ReadyAnnouncements uint `json:"readyAnnouncements"`

	// ReadySessions is the last number of Session objects that are ready (= e.g. Established for BGP) on the actual routing daemon.
	ReadySessions uint `json:"readySessions"`

	// ProgressingAnnouncements is the last number of Announcement objects that are currently progressing on the actual routing daemon.
	ProgressingAnnouncements uint `json:"progressingAnnouncements"`

	// ProgressingSessions is the last number of Session objects that are currently progressing on the actual routing daemon.
	ProgressingSessions uint `json:"progressingSessions"`

	// FailedAnnouncements is the last observed number of Announcement objects that are failed on the actual routing daemon.
	FailedAnnouncements uint `json:"failedAnnouncements"`

	// FailedSessions is the last observed number of Sessions objects that are failed on the actual routing daemon.
	FailedSessions uint `json:"failedSessions"`

	// When was the status last updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`

	// Represents the latest available observations of the status of this Router.
	Conditions []RouterCondition `json:"conditions"`
}

// RouterConditionType defines the type of an observed RouterCondition.
type RouterConditionType string

const (
	// RouterReconciling is showing that the given Router is currently in the process of being reconciled.
	RouterReconciling RouterConditionType = RouterConditionType(kstatus.ConditionReconciling)

	// RouterStalled is showing that the given Router is currently not being able to be reconciled for some reason.
	RouterStalled RouterConditionType = RouterConditionType(kstatus.ConditionStalled)
)

// RouterCondition represents an observed condition of a given Router instance.
type RouterCondition struct {
	CommonCondition `json:",inline"`

	// The type of this Router condition
	Type RouterConditionType `json:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Router is the Schema for the routers API
type Router struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouterSpec   `json:"spec,omitempty"`
	Status RouterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RouterList contains a list of Router
type RouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Router `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Router{}, &RouterList{})
}
