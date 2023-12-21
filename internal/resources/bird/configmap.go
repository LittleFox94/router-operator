package bird

import (
	"context"
	"fmt"
	"net"
	"router-operator/internal/resources"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type templateData struct {
	Router   routerData
	Sessions []sessionData
	Peers    map[string]peerData
}

type routerData struct {
	ID string
}

type sessionData struct {
	Name      string
	Peer      string
	SourceIP  string
	PeerIP    string
	PeerRange bool
	MyASN     uint
}

type peerData struct {
	Name string
	ASN  uint
}

func reconcileConfigMaps(ctx context.Context, c client.Client, reconciliation *resources.RouterReconciliation) error {
	data := templateData{
		Router: routerData{
			ID: reconciliation.Router.Spec.NodeID.ID,
		},
		Sessions: make([]sessionData, 0),
		Peers:    make(map[string]peerData, 0),
	}

	for _, session := range reconciliation.Sessions {
		peerIP := session.Spec.PeerIP
		peerRange := false

		if strings.Contains(peerIP, "/") {
			if _, n, err := net.ParseCIDR(peerIP); err != nil {
				return fmt.Errorf("error parsing peer IP in CIDR notation: %w", err)
			} else {
				peerRange = true
				peerIP = n.String()
			}
		}

		data.Sessions = append(data.Sessions, sessionData{
			Name:      session.GetName(),
			Peer:      session.Spec.Peer.Name,
			SourceIP:  session.Spec.SourceIP,
			PeerIP:    peerIP,
			PeerRange: peerRange,
			MyASN:     session.Spec.BGP.MyASN,
		})
	}

	for _, peer := range reconciliation.Peers {
		data.Peers[peer.Name] = peerData{
			Name: peer.GetName(),
			ASN:  peer.Spec.BGP.ASN,
		}
	}

	content := strings.Builder{}
	if err := configTemplate.Execute(&content, data); err != nil {
		return fmt.Errorf("error rendering template: %w", err)
	}

	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: reconciliation.Router.GetNamespace(),
			Name:      fmt.Sprintf("router-%s-config", reconciliation.Router.GetName()),
		},
	}

	status, err := controllerutil.CreateOrUpdate(ctx, c, &cm, func() error {
		cm.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			*metav1.NewControllerRef(reconciliation.Router, reconciliation.Router.GetObjectKind().GroupVersionKind()),
		}

		for label, value := range reconciliation.Router.ObjectMeta.GetLabels() {
			metav1.SetMetaDataLabel(&cm.ObjectMeta, label, value)
		}

		for annotation, value := range reconciliation.Router.ObjectMeta.GetAnnotations() {
			metav1.SetMetaDataAnnotation(&cm.ObjectMeta, annotation, value)
		}

		if cm.Data == nil {
			cm.Data = make(map[string]string)
		}

		cm.Data["bird.conf"] = content.String()

		return nil
	})

	if err != nil {
		return fmt.Errorf("error creating or updating ConfigMap: %w", err)
	}

	log.FromContext(ctx).Info("ConfigMap reconciled", "ConfigMap", client.ObjectKeyFromObject(&cm), "status", status)
	return nil
}
