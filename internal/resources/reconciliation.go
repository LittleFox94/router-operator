package resources

import (
	routingv1alpha1 "router-operator/api/v1alpha1"
)

type RouterReconciliation struct {
	Router        *routingv1alpha1.Router
	Sessions      []*routingv1alpha1.Session
	Peers         []*routingv1alpha1.Peer
	Announcements []*routingv1alpha1.Announcement
}
