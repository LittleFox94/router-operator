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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PeerSpec defines the desired state of Peer
type PeerSpec struct {
	// BGP peer attributes.
	BGP *BGPPeerSpec `json:"bgp,omitempty"`
}

// BGPPeerSpec contains the BGP-specific attributes of a Peer.
type BGPPeerSpec struct {
	// ASN of the peer.
	ASN uint `json:"asn"`
}

// PeerStatus defines the observed state of Peer
type PeerStatus struct {
	// Generation of this Peer at the time of updating this status.
	ObservedGeneration int64 `json:"observedGeneration"`

	// ClusterSessions is the last number of Session objects observed for this Peer.
	ClusterSessions uint `json:"clusterSessions"`

	// ReadySessions is the last number of Session objects that are ready (= e.g. Established for BGP) on the actual routing daemon.
	ReadySessions uint `json:"readySessions"`

	// ProgressingSessions is the last number of Session objects that are progressing (currently being configured or coming up) on the actual routing daemon.
	ProgressingSessions uint `json:"progressingSessions"`

	// FailedSessions is the last observed number of Sessions objects that are failed on the actual routing daemon.
	FailedSessions uint `json:"failedSessions"`

	// When was the status last updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Peer is the Schema for the peers API
type Peer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PeerSpec   `json:"spec,omitempty"`
	Status PeerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PeerList contains a list of Peer
type PeerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Peer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Peer{}, &PeerList{})
}
