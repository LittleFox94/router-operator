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

// SessionSpec defines the desired state of Session
type SessionSpec struct {
	// Which Router should establish this session.
	Router corev1.LocalObjectReference `json:"router"`

	// Which Peer is on the other side of this session.
	Peer corev1.LocalObjectReference `json:"peer"`

	// IP to listen on and used to connect to the Peer.
	SourceIP string `json:"sourceIP"`

	// IP of the Peer to connect to.
	PeerIP string `json:"peerIP"`

	// BGP session attributes.
	BGP *BGPSessionSpec `json:"bgp,omitempty"`
}

// BGPSessionSpec contains the BGP-specific attributes of a Sessions.
type BGPSessionSpec struct {
	// ASN on my side of the session.
	MyASN uint `json:"myASN"`
}

// SessionStatus defines the observed state of Session
type SessionStatus struct {
	// Generation of this Session at the time of updating this status.
	ObservedGeneration int64 `json:"observedGeneration"`

	// Routes exported via this Session.
	Exported uint `json:"exported,omitempty"`

	// Routes imported via this Session.
	Imported uint `json:"imported,omitempty"`

	// When was the status last updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`

	// Represents the latest available observations of the status of this Session.
	Conditions []SessionCondition `json:"conditions"`
}

// SessionConditionType defines the type of an observed SessionCondition.
type SessionConditionType string

const (
	// SessionReconciling is showing that the given Session is currently in the process of being reconciled.
	SessionReconciling RouterConditionType = RouterConditionType(kstatus.ConditionReconciling)

	// SessionStalled is showing that the given Session is currently not being able to be reconciled for some reason.
	SessionStalled RouterConditionType = RouterConditionType(kstatus.ConditionStalled)
)

// SessionCondition represents an observed condition of a given Session instance.
type SessionCondition struct {
	CommonCondition `json:",inline"`

	// The type of this Session condition
	Type SessionConditionType `json:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Session is the Schema for the sessions API
type Session struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SessionSpec   `json:"spec,omitempty"`
	Status SessionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SessionList contains a list of Session
type SessionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Session `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Session{}, &SessionList{})
}
