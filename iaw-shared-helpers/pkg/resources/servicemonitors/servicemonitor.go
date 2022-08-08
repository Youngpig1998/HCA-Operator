// ------------------------------------------------------ {COPYRIGHT-TOP} ---
// IBM Confidential
// OCO Source Materials
// 5900-AEO
//
// Copyright IBM Corp. 2021
//
// The source code for this program is not published or otherwise
// divested of its trade secrets, irrespective of what has been
// deposited with the U.S. Copyright Office.
// ------------------------------------------------------ {COPYRIGHT-END} ---
package servicemonitors

import (
	"github.com/Youngpig1998/HCA-Operator/iaw-shared-helpers/pkg/resources"
	prometheusv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceMonitor is a balabala
// Reconcileable interface
type ServiceMonitor struct {
	*prometheusv1.ServiceMonitor
}

// From returns a new Reconcileable ServiceMonitor from a prometheusv1.ServiceMonitor
func From(serviceMonitor *prometheusv1.ServiceMonitor) *ServiceMonitor {
	return &ServiceMonitor{ServiceMonitor: serviceMonitor}
}

// ShouldUpdate returns whether the resource should be updated in Kubernetes and
// the resource to update with
func (s ServiceMonitor) ShouldUpdate(current client.Object) (bool, client.Object) {
	currentServiceMonitor := current.DeepCopyObject().(*prometheusv1.ServiceMonitor)
	newServiceMonitor := currentServiceMonitor.DeepCopy()
	resources.MergeMetadata(newServiceMonitor, s)
	newServiceMonitor.Spec = s.Spec
	return !equality.Semantic.DeepEqual(newServiceMonitor, current), newServiceMonitor
}

// GetResource retrieves the resource instance
func (s ServiceMonitor) GetResource() client.Object {
	return s.ServiceMonitor
}

// ResourceKind retrieves the string kind of the resource
func (s ServiceMonitor) ResourceKind() string {
	return "ServiceMonitor"
}

// ResourceIsNil returns whether or not the resource is nil
func (s ServiceMonitor) ResourceIsNil() bool {
	return s.ServiceMonitor == nil
}

// NewResourceInstance returns a new instance of the same resource type
func (s ServiceMonitor) NewResourceInstance() client.Object {
	return &prometheusv1.ServiceMonitor{}
}
