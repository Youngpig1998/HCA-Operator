/*
Copyright 2022.

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

package v1beta1

import (
	prometheusv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HCAJobSpec defines the desired state of HCAJob
type HCAJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value". The requirements are ANDed.
	//MatchLabels map[string]string `json:"matchLabels" protobuf:"bytes,1,rep,name=matchLabels"`

	// updateInterval is an int32 variable. It is used to set time interval of querying the cluster
	//UpdateInterval int32 `json:"updateInterval"`

	// clusterData is a custom type.
	//ClusterData ClusterDataSpec `json:"clusterData"`

	// the namespace where microservice app was deployed in cluster
	AppNamespace string `json:"appNamespace"`

	// scaleData is a custom type. Now, we can just auto-scale the microservice Deployment resource.
	ScaleDatas ScaleDataSpec `json:"scaleDatas" protobuf:"bytes,4,rep,name=scaleDatas"`

	//monitorData is just a wrapper of prometheusv1.ServiceMonitorSpec.
	//The only different is that the NamespaceSelector field cannot be empty.
	MonitorDatas MonitorDataSpec `json:"monitorDatas" protobuf:"bytes,1,rep,name=monitorDatas"`
}

// HCAJobStatus defines the observed state of HCAJob
type HCAJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HCAJob is the Schema for the hcajobs API
type HCAJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCAJobSpec   `json:"spec,omitempty"`
	Status HCAJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HCAJobList contains a list of HCAJob
type HCAJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCAJob `json:"items"`
}

type ClusterDataSpec struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	Properties *runtime.RawExtension `json:"properties,omitempty"`
}

type ScaleDataSpec struct {
	// scaleTargetRef points to the target resource to scale, and is used to the pods for which metrics
	// should be collected, as well as to actually change the replica count.
	//ScaleTargetRef autoscalingv2beta2.CrossVersionObjectReference `json:"scaleTargetRef" protobuf:"bytes,1,opt,name=scaleTargetRef"`
	//Specify the microservice Deployments we want to auto-scale
	ScaleTargetDeploymentNames []string `json:"scaleTargetDeploymentNames"`
	// minReplicas is the lower limit for the number of replicas to which the autoscaler
	// can scale down.  It defaults to 1 pod.  minReplicas is allowed to be 0 if the
	// alpha feature gate HPAScaleToZero is enabled and at least one Object or External
	// metric is configured.  Scaling is active as long as at least one metric value is
	// available.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty" protobuf:"varint,2,opt,name=minReplicas"`
	// maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// It cannot be less that minReplicas.
	MaxReplicas int32 `json:"maxReplicas" protobuf:"varint,3,opt,name=maxReplicas"`
	// metrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).  The desired replica count is calculated multiplying the
	// ratio between the target value and the current value by the current
	// number of pods.  Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source types for
	// more information about how each type of metric must respond.
	// If not set, the default metric will be set to 80% average CPU utilization.
	// +optional
	Metrics []autoscalingv2beta2.MetricSpec `json:"metrics,omitempty" protobuf:"bytes,4,rep,name=metrics"`
	// behavior configures the scaling behavior of the target
	// in both Up and Down directions (scaleUp and scaleDown fields respectively).
	// If not set, the default HPAScalingRules for scale up and scale down are used.
	// +optional
	Behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty" protobuf:"bytes,5,opt,name=behavior"`
}

type MonitorDataSpec struct {
	// The label to use to retrieve the job name from.
	JobLabel string `json:"jobLabel,omitempty"`
	// TargetLabels transfers labels on the Kubernetes Service onto the target.
	TargetLabels []string `json:"targetLabels,omitempty"`
	// PodTargetLabels transfers labels on the Kubernetes Pod onto the target.
	PodTargetLabels []string `json:"podTargetLabels,omitempty"`
	// A list of endpoints allowed as part of this ServiceMonitor.
	Endpoints []prometheusv1.Endpoint `json:"endpoints"`
	// A list of app Services used by ServiceMonitors. The relationship is one to one,
	//which means a Service correponds to a ServiceMonitor
	ServiceLabels []map[string]string `json:"serviceLabels"`
	// Selector to select which namespaces the Endpoints objects are discovered from.
	//NamespaceSelector prometheusv1.NamespaceSelector `json:"namespaceSelector"`
	// SampleLimit defines per-scrape limit on number of scraped samples that will be accepted.
	SampleLimit uint64 `json:"sampleLimit,omitempty"`
}

func init() {
	SchemeBuilder.Register(&HCAJob{}, &HCAJobList{})
}
